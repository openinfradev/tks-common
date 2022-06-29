package argowf

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/openinfradev/tks-common/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apiclient "github.com/argoproj/argo-workflows/v3/pkg/apiclient"
	workflowpkg "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo-workflows/v3/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Client interface {
	GetWorkflow(ctx context.Context, workflowId string, namespace string) (wfv1.WorkflowPhase, string, error)
	SumbitWorkflowFromWftpl(ctx context.Context, wftplName string, namespace string, parameters []string) (string, error)
	IsRunningWorkflowByContractId(ctx context.Context, nameSpace string, contractId string) (bool, error)
	WaitWorkflows(ctx context.Context, namespace string, workflowNames []string, ignoreNotFound, quiet bool) bool
}

type ArgoClient struct {
	serviceClient workflowpkg.WorkflowServiceClient
}

func New(host string, port int) (Client, error) {
	baseUrl := fmt.Sprintf("%s:%d", host, port)
	opts := apiclient.Opts{
		ArgoServerOpts: apiclient.ArgoServerOpts{
			URL:                baseUrl,
			Secure:             false,
			InsecureSkipVerify: false,
			HTTP1:              false,
		},
		AuthSupplier: func() string {
			return ""
		},
	}
	_, client, err := apiclient.NewClientFromOpts(opts)
	if err != nil {
		return nil, err
	}
	serviceClient := client.NewWorkflowServiceClient()

	return &ArgoClient{
		serviceClient: serviceClient,
	}, nil
}

func (c *ArgoClient) SumbitWorkflowFromWftpl(ctx context.Context, wftplName string, namespace string, parameters []string) (string, error) {
	submitOpts := wfv1.SubmitOpts{}
	submitOpts.Parameters = parameters

	created, err := c.serviceClient.SubmitWorkflow(ctx, &workflowpkg.WorkflowSubmitRequest{
		Namespace:     namespace,
		ResourceKind:  "WorkflowTemplate",
		ResourceName:  wftplName,
		SubmitOptions: &submitOpts,
	})

	if err != nil {
		log.Error("Failed to submit argo workflow. err ", err)
		return "", err
	}

	return created.Name, nil
}

func (c *ArgoClient) GetWorkflow(ctx context.Context, workflowId string, namespace string) (wfv1.WorkflowPhase, string, error) {
	workflow, err := c.serviceClient.GetWorkflow(ctx, &workflowpkg.WorkflowGetRequest{
		Name:      workflowId,
		Namespace: "argo",
	})
	if err != nil {
		log.Error("Failed to get argo workflow. err ", err)
		return "", "", err
	}

	message := fmt.Sprintf("(%s) %s", workflow.Status.Progress, workflow.Status.Message)

	return workflow.Status.Phase, message, nil
}

func (c *ArgoClient) IsRunningWorkflowByContractId(ctx context.Context, nameSpace string, contractId string) (bool, error) {
	wfList, err := c.serviceClient.ListWorkflows(ctx, &workflowpkg.WorkflowListRequest{
		Namespace:   nameSpace,
		ListOptions: &metav1.ListOptions{LabelSelector: "workflows.argoproj.io/phase in (Pending,Running)"},
		Fields:      "items.metadata.name,items.spec",
	})

	if err != nil {
		log.Error("failed to get argo workflows. err : ", err)
		return false, err
	}

	for _, item := range wfList.Items {
		for _, arg := range item.Spec.Arguments.Parameters {
			if arg.Name == "contract_id" && arg.Value.String() == contractId {
				return true, errors.New("Existed running(pending) workflow")
			}
		}
	}

	return false, nil
}

func (c *ArgoClient) WaitWorkflows(ctx context.Context, namespace string, workflowNames []string, ignoreNotFound, quiet bool) bool {
	log.Debug("waiting workflowNames : ", workflowNames)

	var wg sync.WaitGroup
	wfSuccessStatus := true

	for _, name := range workflowNames {
		wg.Add(1)
		go func(name string) {
			if !c.waitOnOne(ctx, name, namespace, ignoreNotFound, quiet) {
				wfSuccessStatus = false
			}
			wg.Done()
		}(name)

	}
	wg.Wait()

	return wfSuccessStatus
}

func (c *ArgoClient) waitOnOne(ctx context.Context, wfName, namespace string, ignoreNotFound, quiet bool) bool {
	req := &workflowpkg.WatchWorkflowsRequest{
		Namespace: namespace,
		ListOptions: &metav1.ListOptions{
			FieldSelector:   util.GenerateFieldSelectorFromWorkflowName(wfName),
			ResourceVersion: "0",
		},
	}
	stream, err := c.serviceClient.WatchWorkflows(ctx, req)
	if err != nil {
		if status.Code(err) == codes.NotFound && ignoreNotFound {
			return true
		}
		return false
	}
	for {
		event, err := stream.Recv()
		if err == io.EOF {
			stream, _ = c.serviceClient.WatchWorkflows(ctx, req)
			continue
		}
		if event == nil {
			continue
		}
		wf := event.Object
		if !wf.Status.FinishedAt.IsZero() {
			if !quiet {
				log.Info(fmt.Sprintf("%s %s at %v\n", wfName, wf.Status.Phase, wf.Status.FinishedAt))
			}
			if wf.Status.Phase == wfv1.WorkflowFailed || wf.Status.Phase == wfv1.WorkflowError {
				return false
			}
			return true
		}
	}
}
