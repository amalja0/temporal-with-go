package temporal

import (
	"fmt"
	"sales-record-orchestration/internal/domain"
	"sales-record-orchestration/internal/ports"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type WorkflowWorker interface {
	SalesETLWorkflow(ctx workflow.Context, params domain.SalesQueryParams) error
}

type workflowWorker struct {
	activity ports.TemporalActivity
}

func InitWorkflowWorker(activity ports.TemporalActivity) WorkflowWorker {
	return &workflowWorker{
		activity: activity,
	}
}

func (w *workflowWorker) SalesETLWorkflow(ctx workflow.Context, params domain.SalesQueryParams) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var sales []domain.SaleWithDetail
	if err := workflow.ExecuteActivity(ctx, w.activity.FetchSalesActivity, params).Get(ctx, &sales); err != nil {
		fmt.Println("error on first activity", err)
		return err
	}

	if err := workflow.ExecuteActivity(ctx, w.activity.PublishSalesActivity, sales).Get(ctx, nil); err != nil {
		fmt.Println(err)
		return err
	}

	var orderRecords []domain.OrderRecord
	for _, sale := range sales {
		orderRecords = append(orderRecords, domain.OrderRecord{
			Id:          uuid.New(),
			SaleId:      sale.ID,
			Quantity:    sale.Qty,
			SaleAmount:  sale.SaleAmount,
			Discount:    sale.Discount,
			Profit:      sale.Profit,
			ProfitRatio: sale.ProfitRatio,
			OrderId:     sale.OrderID,
			OrderDate:   sale.OrderDate,
			LocationId:  sale.LocationID,
			ProductId:   sale.ProductID,
			SegmentId:   sale.SegmentID,
			ProductName: sale.ProductName,
			SegmentName: sale.SegmentName,
			CreatedAt:   time.Now(),
		})
	}

	if err := workflow.ExecuteActivity(ctx, w.activity.ProcessSalesActivity, orderRecords).Get(ctx, &orderRecords); err != nil {
		fmt.Println("ERRROROR", err)
		return err
	}

	return nil
}
