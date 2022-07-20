package test

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"time"
)

const (
	Region            = "us-east-1"
	RetryDelaySeconds = 30
	RetryAttempts     = 20
)

type GrafanaQuery struct {
	IntervalMs    int64  `json:"intervalMs"`
	MaxDataPoints int    `json:"maxDataPoints"`
	DatasourceId  int    `json:"datasourceId"`
	RawSql        string `json:"rawSql"`
	Format        string `json:"format"`
}

type GrafanaQueryRequest struct {
	From    string          `json:"from"`
	To      string          `json:"to"`
	Queries []*GrafanaQuery `json:"queries"`
}

// getAWSSession Logs in to AWS and return a session
func getAWSSession() *session.Session {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		panic(err)
	}
	return sess
}

// countRecords gets a count of records from the table and compares it to the count expected
func countRecords(awsSession *session.Session, dbName string, dbArn string, dbSecretsArn string, table string, owner string, repo string, minExpected int) error {
	for i := 0; ; i++ {
		fmt.Printf("Getting count of rows in %s table for %s/%s\n", table, owner, repo)
		svc := rdsdataservice.New(awsSession, aws.NewConfig().WithRegion(Region))
		output, err := svc.ExecuteStatement(&rdsdataservice.ExecuteStatementInput{
			Database:    aws.String(dbName),
			ResourceArn: aws.String(dbArn),
			SecretArn:   aws.String(dbSecretsArn),
			Sql:         aws.String(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE owner = '%s' AND repo = '%s'", table, owner, repo)),
		})
		if err != nil {
			fmt.Println(err)
		} else {
			count := int(*output.Records[0][0].LongValue)
			fmt.Println("Count: ", count)
			if count >= minExpected {
				return nil
			}
		}

		if i >= (RetryAttempts - 1) {
			panic("Timed out while retrying")
		}

		fmt.Printf("Retrying in %d seconds...\n", RetryDelaySeconds)
		time.Sleep(time.Second * RetryDelaySeconds)
	}
}

// checkForEmptyFields checks for any records which have empty fields
func checkForEmptyFields(awsSession *session.Session, dbName string, dbArn string, dbSecretsArn string, table string) error {
	fmt.Printf("Checking for count of records with empty fields in %s table\n", table)
	svc := rdsdataservice.New(awsSession, aws.NewConfig().WithRegion(Region))

	query := fmt.Sprintf(`SELECT count(*) FROM %s
                                 WHERE owner IS NULL OR owner = ''
                                 OR repo IS NULL OR repo = ''
                                 OR branch IS NULL OR branch = ''`, table)

	output, err := svc.ExecuteStatement(&rdsdataservice.ExecuteStatementInput{
		Database:    aws.String(dbName),
		ResourceArn: aws.String(dbArn),
		SecretArn:   aws.String(dbSecretsArn),
		Sql:         aws.String(query),
	})
	if err != nil {
		return err
	}
	count := int(*output.Records[0][0].LongValue)
	fmt.Println("Count: ", count)
	if count > 0 {
		return fmt.Errorf("found %d records with empty fields", count)
	}
	return nil
}

// dropTable drops a table from the database
func dropTable(awsSession *session.Session, dbName string, dbArn string, dbSecretsArn string, table string) error {
	fmt.Printf("Dropping table: %s\n", table)
	svc := rdsdataservice.New(awsSession, aws.NewConfig().WithRegion(Region))

	query := fmt.Sprintf(`drop table if exists %s`, table)

	_, err := svc.ExecuteStatement(&rdsdataservice.ExecuteStatementInput{
		Database:    aws.String(dbName),
		ResourceArn: aws.String(dbArn),
		SecretArn:   aws.String(dbSecretsArn),
		Sql:         aws.String(query),
	})
	if err != nil {
		return err
	}
	return nil
}

func validateReruns(awsSession *session.Session, dbName string, dbArn string, dbSecretsArn string, table string, owner string, repo string, minExpected int) error {
	fmt.Printf("Checking for count of records with empty fields in %s table\n", table)
	svc := rdsdataservice.New(awsSession, aws.NewConfig().WithRegion(Region))

	query := fmt.Sprintf(`SELECT count(*) FROM %s
                                 WHERE node_id IN (SELECT node_id FROM %s WHERE run_attempt > 1)
                                 AND owner = "%s" AND repo = "%s"`, table, table, owner, repo)

	output, err := svc.ExecuteStatement(&rdsdataservice.ExecuteStatementInput{
		Database:    aws.String(dbName),
		ResourceArn: aws.String(dbArn),
		SecretArn:   aws.String(dbSecretsArn),
		Sql:         aws.String(query),
	})
	if err != nil {
		return err
	}
	count := int(*output.Records[0][0].LongValue)
	fmt.Println("Count: ", count)
	if count < minExpected {
		return fmt.Errorf("found %d rerun records. expected %d", count, minExpected)
	}
	return nil
}

// deleteRecentCommits deletes a specific number of records from the table
func deleteRecentCommits(awsSession *session.Session, dbName string, dbArn string, dbSecretsArn string, table string, owner string, repo string, rows int) error {
	fmt.Printf("Deleting most recent %d rows in %s table for %s/%s\n", rows, table, owner, repo)
	svc := rdsdataservice.New(awsSession, aws.NewConfig().WithRegion(Region))
	_, err := svc.ExecuteStatement(&rdsdataservice.ExecuteStatementInput{
		Database:    aws.String(dbName),
		ResourceArn: aws.String(dbArn),
		SecretArn:   aws.String(dbSecretsArn),
		Sql:         aws.String(fmt.Sprintf("DELETE FROM %s WHERE OWNER = '%s' AND repo = '%s' ORDER BY committer_date DESC LIMIT %d;", table, owner, repo, rows)),
	})
	if err != nil {
		return err
	}
	return nil
}
