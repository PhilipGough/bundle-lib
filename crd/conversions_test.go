package crd

import (
	"github.com/automationbroker/broker-client-go/pkg/apis/automationbroker.io/v1"
	"github.com/automationbroker/bundle-lib/apb"
	"reflect"
	"testing"
)

func TestConvertSpecToBundle(t *testing.T) {
	mockSpec := getMockSpec()
	result, err := ConvertSpecToBundle(mockSpec)
	expect := getMockBundle()
	if err != nil {
		t.Errorf("unexpected error %s", err.Error())
	}

	if !reflect.DeepEqual(result, expect) {
		t.Errorf("unexpected results, expected %+v but got \n %+v", expect, result)
	}
}

func TestConvertBundleToSpec(t *testing.T) {
	mockBundle := getMockBundle()
	result, err := ConvertBundleToSpec(mockBundle, "Test")
	expect := getMockSpec()
	if err != nil {
		t.Errorf("unexpected error %s", err.Error())
	}

	if !reflect.DeepEqual(result, expect) {
		t.Errorf("unexpected results, expected %+v but got /n %+v", expect, result)
	}
}

func TestConvertServiceInstanceToCRD(t *testing.T) {
	mockSi := getMockServiceInstance()
	result, err := ConvertServiceInstanceToCRD(mockSi)
	expect := getMockServiceInstanceSpec()
	if err != nil {
		t.Errorf("unexpected error %s", err.Error())
	}

	if !reflect.DeepEqual(result, expect) {
		t.Errorf("unexpected results, expected %+v but got %+v", expect, result)
	}
}

func TestConvertServiceInstanceToAPB(t *testing.T) {
	siSpec := getMockServiceInstanceSpec()
	spec := getMockSpec()
	result, err := ConvertServiceInstanceToAPB(siSpec, spec, "f4b72d2a-bf85-4847-a19f-edd8beef0ddd")
	expect := getMockServiceInstance()
	if err != nil {
		t.Errorf("unexpected error %s", err.Error())
	}

	if !reflect.DeepEqual(result, expect) {
		t.Errorf("unexpected results, expected %+v but got %+v", expect, result)
	}
}

func TestConvertStateToCRD(t *testing.T) {
	inputs := getStateInputs()
	const FakeState apb.State = "invalid"
	inputs = append(inputs, stateConversions{FakeState, v1.StateFailed})
	for _, input := range inputs {
		result := ConvertStateToCRD(input.apbState)
		if result != input.brokerState {
			t.Errorf("unexpected result. expect %s got $s", input.brokerState, result)
		}
	}
}

func TestConvertStateToAPB(t *testing.T) {
	inputs := getStateInputs()
	const FakeState v1.State = "invalid"
	inputs = append(inputs, stateConversions{apb.StateFailed, FakeState})
	for _, input := range inputs {
		result := ConvertStateToAPB(input.brokerState)
		if result != input.apbState {
			t.Errorf("unexpected result. expect %s got $s", input.apbState, result)
		}
	}
}

func TestConvertJobMethodToAPB(t *testing.T) {
	inputs := getJobMethodInputs()
	const FakeJobMetod v1.JobMethod = "invalid"
	inputs = append(inputs, jobMethodConversion{FakeJobMetod, apb.JobMethodProvision})
	for _, input := range inputs {
		result := convertJobMethodToAPB(input.brokerJobMethod)
		if result != input.apbJobMethod {
			t.Errorf("unexpected result, expect %s got %s", input.apbJobMethod, result)
		}
	}
}

func TestConvertJobMethodToCRD(t *testing.T) {
	inputs := getJobMethodInputs()
	const FakeJobMethod apb.JobMethod = "invalid"
	inputs = append(inputs, jobMethodConversion{v1.JobMethodProvision, FakeJobMethod})
	for _, input := range inputs {
		result := convertJobMethodToCRD(input.apbJobMethod)
		if result != input.brokerJobMethod {
			t.Errorf("unexpected result, expect %s got %s", input.brokerJobMethod, result)
		}
	}
}

func TestConvertAsyncTypeToString(t *testing.T) {
	inputs := getAsyncToStringInputs()
	const FakeAsync v1.AsyncType = "invalid"
	inputs = append(inputs, asyncConversion{FakeAsync, "required"})
	for _, input := range inputs {
		result := convertAsyncTypeToString(input.brokerConst)
		if result != input.stringVal {
			t.Errorf("unexpected result, expect %s go %s", input.stringVal, result)
		}
	}
}

func TestConvertToAsyncType(t *testing.T) {
	inputs := getAsyncToStringInputs()
	inputs = append(inputs, asyncConversion{v1.RequiredAsync, "anyVal"})
	for _, input := range inputs {
		result := convertToAsyncType(input.stringVal)
		if result != input.brokerConst {
			t.Errorf("unexpected result, expect %s go %s", input.stringVal, result)
		}
	}
}

type asyncConversion struct {
	brokerConst v1.AsyncType
	stringVal   string
}

type jobMethodConversion struct {
	brokerJobMethod v1.JobMethod
	apbJobMethod    apb.JobMethod
}

type stateConversions struct {
	apbState    apb.State
	brokerState v1.State
}

func getMockSpec() *apb.Spec {
	return &apb.Spec{
		ID:          "Test",
		Runtime:     1,
		Version:     "1",
		FQName:      "MockKeycloak",
		Image:       "test",
		Tags:        []string{"auth service"},
		Bindable:    true,
		Description: "Fake Keycloak service",
		Metadata: map[string]interface{}{
			"displayName":      "MockKeyCloak",
			"documentationUrl": "http://www.keycloak.org/documentation.html",
			"dependencies": []string{
				"docker.io/jboss/keycloak-openshift:3.4.3.Final",
				"docker.io/centos/postgresql-96-centos7:9.6",
			},
		},
		Async: "required",
		Plans: []apb.Plan{
			{ID: "Mock", Name: "MockPlan", Description: "Mock Plan"},
		},
	}
}

func getMockBundle() v1.BundleSpec {
	return v1.BundleSpec{
		Runtime:     1,
		Version:     "1",
		FQName:      "MockKeycloak",
		Image:       "test",
		Tags:        []string{"auth service"},
		Bindable:    true,
		Async:       "required",
		Description: "Fake Keycloak service",
		Metadata:    "{\"dependencies\":[\"docker.io/jboss/keycloak-openshift:3.4.3.Final\",\"docker.io/centos/postgresql-96-centos7:9.6\"],\"displayName\":\"MockKeyCloak\",\"documentationUrl\":\"http://www.keycloak.org/documentation.html\"}",
		Plans:       []v1.Plan{{"Mock", "MockPlan", "Mock Plan", "null", false, false, []string{}, []v1.Parameters{}, []v1.Parameters{}}},
	}
}

func getMockServiceInstance() *apb.ServiceInstance {
	bindings := make(map[string]bool)
	bindings["one"] = true
	bindings["two"] = false

	params := make(map[string]interface{})
	params["first"] = "1"
	params["second"] = "2"

	return &apb.ServiceInstance{
		ID:   []byte("f4b72d2a-bf85-4847-a19f-edd8beef0ddd"),
		Spec: getMockSpec(),
		Context: &apb.Context{
			Platform:  "kubernetes",
			Namespace: "test",
		},
		Parameters: &apb.Parameters{
			"key": "value",
		},
		BindingIDs: bindings,
	}
}

func getMockServiceInstanceSpec() v1.ServiceInstanceSpec {
	return v1.ServiceInstanceSpec{
		BundleID: "Test",
		Context: v1.Context{
			Plateform: "kubernetes",
			Namespace: "test",
		},
		Parameters: "{\"key\" : \"value\"}",
		BindingIDs: []string{"one", "two"},
	}
}

func getStateInputs() []stateConversions {
	return []stateConversions{
		{
			apb.StateNotYetStarted,
			v1.StateNotYetStarted,
		},
		{
			apb.StateInProgress,
			v1.StateInProgress,
		},
		{
			apb.StateSucceeded,
			v1.StateSucceeded,
		},
		{
			apb.StateFailed,
			v1.StateFailed,
		},
	}
}

func getJobMethodInputs() []jobMethodConversion {
	return []jobMethodConversion{
		{
			v1.JobMethodProvision,
			apb.JobMethodProvision,
		},
		{
			v1.JobMethodDeprovision,
			apb.JobMethodDeprovision,
		},
		{
			v1.JobMethodBind,
			apb.JobMethodBind,
		},
		{
			v1.JobMethodUnbind,
			apb.JobMethodUnbind,
		},
		{
			v1.JobMethodUpdate,
			apb.JobMethodUpdate,
		},
	}
}

func getAsyncToStringInputs() []asyncConversion {
	return []asyncConversion{
		{
			v1.OptionalAsync,
			"optional",
		},
		{
			v1.RequiredAsync,
			"required",
		},
		{
			v1.Unsupported,
			"unsupported",
		},
	}
}
