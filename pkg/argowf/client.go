package argowf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/openinfradev/tks-common/pkg/log"
)

type Client interface {
	GetWorkflowTemplates(namespace string) (*GetWorkflowTemplatesResponse, error)
	GetWorkflow(namespace string, workflowName string) (*Workflow, error)
	GetWorkflows(namespace string) (*GetWorkflowsResponse, error)
	SumbitWorkflowFromWftpl(wftplName, targetNamespace string, opts SubmitOptions) (string, error)
}

type ArgoClient struct {
	client *http.Client
	url    string
}

// New
func New(host string, port int, ssl bool, token string) (Client, error) {
	var baseUrl string
	if ssl {
		if token == "" {
			return nil, fmt.Errorf("argo ssl enabled but token is empty.")
		}
		baseUrl = fmt.Sprintf("https://%s:%d", host, port)
	} else {
		baseUrl = fmt.Sprintf("http://%s:%d", host, port)
	}
	return &ArgoClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns: 10,
			},
		},
		url: baseUrl,
	}, nil
}

func (c *ArgoClient) GetWorkflowTemplates(namespace string) (*GetWorkflowTemplatesResponse, error) {
	res, err := http.Get(fmt.Sprintf("%s/api/v1/workflow-templates/%s", c.url, namespace))
	if err != nil || res == nil {
		log.Error("error from get workflow-templats err: ", err)
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Error("error from get workflow-templats return code: ", res.StatusCode)
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Error("error closing http body")
		}
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	wftplRes := GetWorkflowTemplatesResponse{}
	if err := json.Unmarshal(body, &wftplRes); err != nil {
		log.Error("an error was unexpected while parsing response from api /workflow template.")
		return nil, err
	}
	return &wftplRes, nil
}

func (c *ArgoClient) GetWorkflow(namespace string, workflowName string) (*Workflow, error) {
	res, err := http.Get(fmt.Sprintf("%s/api/v1/workflows/%s/%s", c.url, namespace, workflowName))
	if err != nil || res == nil {
		log.Error("error from get workflow err: ", err)
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Error("error from get workflow return code: ", res.StatusCode)
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Error("error closing http body")
		}
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	workflowRes := Workflow{}
	if err := json.Unmarshal(body, &workflowRes); err != nil {
		log.Error("an error was unexpected while parsing response from api /workflow template.")
		return nil, err
	}

	return &workflowRes, nil
}

func (c *ArgoClient) GetWorkflows(namespace string) (*GetWorkflowsResponse, error) {
	res, err := http.Get(fmt.Sprintf("%s/api/v1/workflows/%s", c.url, namespace))
	if err != nil || res == nil {
		log.Error("error from get workflow-templats err: ", err)
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Error("error from get workflow-templats return code: ", res.StatusCode)
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Error("error closing http body")
		}
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	workflowsRes := GetWorkflowsResponse{}
	if err := json.Unmarshal(body, &workflowsRes); err != nil {
		log.Error("an error was unexpected while parsing response from api /workflow template.")
		return nil, err
	}

	return &workflowsRes, nil
}

func (c *ArgoClient) SumbitWorkflowFromWftpl(wftplName, targetNamespace string, opts SubmitOptions) (string, error) {
	reqBody := submitWorkflowRequestBody{
		Namespace:     targetNamespace,
		ResourceKind:  "WorkflowTemplate",
		ResourceName:  wftplName,
		SubmitOptions: opts,
	}
	log.Debug("SumbitWorkflowFromWftpl reqBody ", reqBody)

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "",
			fmt.Errorf("an error was unexpected while marshaling request body")
	}
	buff := bytes.NewBuffer(reqBodyBytes)

	res, err := http.Post(fmt.Sprintf("%s/api/v1/workflows/%s/submit", c.url, targetNamespace), "application/json", buff)

	// [TODO] timeout 처리
	if err != nil || res == nil {
		log.Error("error message ", err.Error())
		return "", err
	}
	if res.StatusCode != 200 {
		log.Error("error from post workflow. return code: ", res.StatusCode)
		return "", err
	}

	defer func() {
		if res != nil {
			if err := res.Body.Close(); err != nil {
				log.Error("error closing http body")
			}
		}
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	submitRes := SubmitWorkflowResponse{}
	if err := json.Unmarshal(body, &submitRes); err != nil {
		log.Error("an error was unexpected while parsing response from api /submit.")
		return "", err
	}
	return submitRes.Metadata.Name, nil
}
