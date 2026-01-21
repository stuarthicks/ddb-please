package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var prettyPrint bool

func fatal(format string, v ...any) {
	fmt.Printf(format+"\n", v...)
	os.Exit(1)
}

func main() {
	var region string
	flag.StringVar(&region, "region", "us-east-1", "AWS Region")

	var list bool
	flag.BoolVar(&list, "list", false, "List tables")

	var info bool
	flag.BoolVar(&info, "info", false, "Get info about a table (e.g. secondary indexes)")

	var get bool
	flag.BoolVar(&get, "get", false, "Get an item")

	var tableName string
	flag.StringVar(&tableName, "t", "", "Name of table to query")

	// var index string
	// flag.StringVar(&index, "i", "", "Name of secondary index to query")

	var partitionKey string
	flag.StringVar(&partitionKey, "p", "", "Partition (HASH) Key")

	var sortKey string
	flag.StringVar(&sortKey, "s", "", "Sort (RANGE) Key")

	flag.BoolVar(&prettyPrint, "pretty", false, "Pretty-print JSON")

	flag.Parse()

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		fatal("unable to load SDK config, %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	if list {
		listTables(ctx, svc)
		return
	}

	if info {
		tableInfo(ctx, svc, tableName)
		return
	}

	if get {
		getItem(ctx, svc, tableName, partitionKey, sortKey)
		return
	}

	// Add other operations here
	// TODO:
	// - local secondary index (must use query)
	// - global secondary index (must use query)
	// - normal queries/scans?

	fatal("select a task using -list/-get, etc. See -h for all available flags.")
}

func listTables(ctx context.Context, svc *dynamodb.Client) {
	in := &dynamodb.ListTablesInput{}
	out, err := svc.ListTables(ctx, in)
	if err != nil {
		fatal("failed to list tables: %s", err.Error())
	}
	for _, tableName := range out.TableNames {
		fmt.Println(tableName)
	}
}

func tableInfo(ctx context.Context, svc *dynamodb.Client, tableName string) {
	in := &dynamodb.DescribeTableInput{TableName: &tableName}
	out, err := svc.DescribeTable(ctx, in)
	if err != nil {
		fatal("failed to describe table: %s", err.Error())
	}

	fmt.Println("Key Schema:")
	for _, a := range out.Table.KeySchema {
		fmt.Printf("\t%s (%s)\n", *a.AttributeName, a.KeyType)
	}

	if len(out.Table.LocalSecondaryIndexes) > 0 {
		fmt.Println()
		fmt.Println("Local Secondary Indexes:")
		for _, lsi := range out.Table.LocalSecondaryIndexes {
			fmt.Printf("\t%s\n\t\t- ", *lsi.IndexName)
			var keys []string
			for _, a := range lsi.KeySchema {
				keys = append(keys, fmt.Sprintf("%s (%s)", *a.AttributeName, a.KeyType))
			}
			fmt.Println(strings.Join(keys, ", "))
		}
	}

	if len(out.Table.GlobalSecondaryIndexes) > 0 {
		fmt.Println()
		fmt.Println("Global Secondary Indexes:")
		for _, gsi := range out.Table.GlobalSecondaryIndexes {
			fmt.Printf("\t%s\n\t\t- ", *gsi.IndexName)
			var keys []string
			for _, a := range gsi.KeySchema {
				keys = append(keys, fmt.Sprintf("%s (%s)", *a.AttributeName, a.KeyType))
			}
			fmt.Println(strings.Join(keys, ", "))
		}
	}
}

func getItem(ctx context.Context, svc *dynamodb.Client, tableName, partitionKey, sortKey string) {
	inDesc := &dynamodb.DescribeTableInput{
		TableName: &tableName,
	}
	outDesc, err := svc.DescribeTable(ctx, inDesc)
	if err != nil {
		fatal("failed to describe table: %s", err.Error())
	}

	key := make(map[string]types.AttributeValue)
	for _, a := range outDesc.Table.KeySchema {
		switch a.KeyType {
		case types.KeyTypeHash:
			key[*a.AttributeName] = &types.AttributeValueMemberS{Value: partitionKey}
		case types.KeyTypeRange:
			key[*a.AttributeName] = &types.AttributeValueMemberS{Value: sortKey}
		}
	}

	in := &dynamodb.GetItemInput{
		TableName: &tableName,
		Key:       key,
	}
	out, err := svc.GetItem(ctx, in)
	if err != nil {
		fatal("failed to get item: %s", err.Error())
	}

	var item map[string]any
	if err := attributevalue.UnmarshalMap(out.Item, &item); err != nil {
		fatal("failed to unmarshal item: %s", err.Error())
	}

	enc := json.NewEncoder(os.Stdout)
	if prettyPrint {
		enc.SetIndent("", "  ")
	}

	if err := enc.Encode(item); err != nil {
		fatal("failed to write item as json to stdout: %s", err.Error())
	}
}
