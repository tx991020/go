package main

import (
	"fmt"

	"github.com/Shopify/sarama"

	"log"

	"os"

	"strings"

	"sync"
)

var (
	wg sync.WaitGroup

	logger = log.New(os.Stderr, "[haha]", log.LstdFlags)
)

func main() {

	sarama.Logger = logger

	consumer, err := sarama.NewConsumer(strings.Split("localhost:9092", ","), nil)

	if err != nil {

		logger.Println("Failed to start consumer: %s", err)

	}

	partitionList, err := consumer.Partitions("hello")

	if err != nil {

		logger.Println("Failed to get the list of partitions: ", err)

	}

	for partition := range partitionList {

		pc, err := consumer.ConsumePartition("hello", int32(partition), sarama.OffsetNewest)

		if err != nil {

			logger.Printf("Failed to start consumer for partition %d: %s\n", partition, err)

		}

		defer pc.AsyncClose()

		wg.Add(1)

		go func(sarama.PartitionConsumer) {

			defer wg.Done()

			for msg := range pc.Messages() {

				fmt.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))

				fmt.Println()

			}

		}(pc)

	}

	wg.Wait()

	logger.Println("Done consuming topic hello")

	consumer.Close()

}
