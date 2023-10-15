package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"peng.fyi/simple-loadtest/runners"
	"peng.fyi/simple-loadtest/tasks"
)

func main() {
	var logger *zap.Logger

	// Step 1. Load YAML configuration file for load runner. See the `example_config.yaml`
	verboseLogging := flag.Bool("v", false, "Enable verbose logging")
	configFilePath := flag.String("config", "", "Path to the YAML configuration file")

	flag.Parse()

	if *verboseLogging {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	if *configFilePath == "" {
		logger.Fatal("provide a YAML configuration file using -config flag")
	}

	c, err := loadConfiguration(*configFilePath)
	if err != nil {
		logger.Fatal("error loading configuration", zap.Error(err))
	}

	// Step 2. Use the config to create Task, TaskRunner and LoadRunner.
	createTaskFunc := func(m map[string]interface{}) func() (tasks.Task, error) {
		taskType := m["type"]
		if taskType == "foo" {
			b, _ := json.Marshal(m)
			var c FooTaskConfig
			if err := json.Unmarshal(b, &c); err != nil {
				return nil
			}
			return func() (tasks.Task, error) {
				return tasks.NewFoo(logger, c.Message), nil
			}
		}
		return nil
	}

	createRunnerFunc := func(c RunnerConfig) func() (*runners.TaskRunner, error) {
		return func() (*runners.TaskRunner, error) {
			r := runners.NewTaskRunner(logger,
				time.Duration(c.LifetimeInSeconds)*time.Second,
				runners.WithIntervalAndJitter(
					time.Duration(c.CycleIntervalInSeconds)*time.Second,
					time.Duration(c.JitterInSeconds)*time.Second))
			return r, nil
		}
	}

	var newRunners []func() (*runners.TaskRunner, error)
	var newTasks []func() (tasks.Task, error)
	for _, c := range c.TaskRunnerConfigs {
		nextRunner := createRunnerFunc(c.Runner)
		nextTask := createTaskFunc(c.Task)
		if nextTask == nil {
			b, _ := json.Marshal(c.Task)
			logger.Error("cannot create task from config", zap.String("config yaml", string(b)))
			continue
		}
		newRunners = append(newRunners, nextRunner)
		newTasks = append(newTasks, nextTask)
	}
	if len(newRunners) == 0 {
		logger.Error("empty runners")
		return
	}
	provider := runners.NewTaskRunnerRR(newRunners, newTasks)
	r := runners.NewLoadRunner(logger, provider,
		time.Duration(c.LifetimeInSeconds)*time.Second,
		c.ConcurrentRunners)
	if err != nil {
		logger.Error("cannot create load")
	}

	// Step 3. Start the load runner and graceful shutdown on signals.
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		r.Start(ctx)
	}()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Info("stop loadtest")
		cancel()
	}()

	wg.Wait()
}

func loadConfiguration(filePath string) (LoadRunnerConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return LoadRunnerConfig{}, err
	}

	var c LoadRunnerConfig

	err = yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		return LoadRunnerConfig{}, err
	}

	return c, nil
}

type LoadRunnerConfig struct {
	LifetimeInSeconds int                `yaml:"lifetimeInSeconds"`
	ConcurrentRunners int                `yaml:"concurrentRunners"`
	TaskRunnerConfigs []TaskRunnerConfig `yaml:"taskRunners"`
}

type TaskRunnerConfig struct {
	Runner RunnerConfig           `yaml:"runner"`
	Task   map[string]interface{} `yaml:"task"`
}

type RunnerConfig struct {
	LifetimeInSeconds      int `yaml:"lifetimeInSeconds"`
	CycleIntervalInSeconds int `yaml:"cycleIntervalInSeconds"`
	JitterInSeconds        int `yaml:"jitterInSeconds"`
}

type FooTaskConfig struct {
	Message string `yaml:"message"`
}
