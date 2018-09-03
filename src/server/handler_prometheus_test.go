package server

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"

	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"github.com/stretchr/testify/suite"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type PrometheusTestSuite struct {
	suite.Suite
	wc *restful.Container
	sp *serviceprovider.Container
}

func (suite *PrometheusTestSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.New(cf)

	//init restful container
	suite.wc = restful.NewContainer()
	suite.sp = serviceprovider.New(cf)
	service := newMonitoringService(sp)
	suite.wc.Add(service)
}

func (suite *PrometheusTestSuite) TearDownSuite() {}

func TestPrometheusSuite(t *testing.T) {
	suite.Run(t, new(PrometheusTestSuite))
}

func (suite *PrometheusTestSuite) TestListNodeMetrics() {
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/nodes/", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestGetNodeMetrics() {
	nodes, err := suite.sp.KubeCtl.GetNodes()
	suite.NoError(err)
	nodeName := nodes[0].GetName()

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/nodes/"+nodeName, nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestListNodeNicsMetrics() {
	nodes, err := suite.sp.KubeCtl.GetNodes()
	suite.NoError(err)
	nodeName := nodes[0].GetName()

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/nodes/"+nodeName+"/nics", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestListPodMetrics() {
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/pods/", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/monitoring/pods?node=.*&namespace=.*&controller=.*", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestGetPodMetrics() {
	namespace := "vortex"
	pods, err := suite.sp.KubeCtl.GetPods(namespace)
	suite.NoError(err)
	podName := pods[0].GetName()

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/pods/"+podName, nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestGetContainerMetrics() {
	namespace := "vortex"
	pods, err := suite.sp.KubeCtl.GetPods(namespace)
	suite.NoError(err)
	podName := pods[0].Name
	containerName := pods[0].Status.ContainerStatuses[0].Name

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/pods/"+podName+"/"+containerName, nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestListServiceMetrics() {
	httpRequest, err := http.NewRequest("GET", `http://172.17.8.100:30003/api/v1/query?query=kube_service_info`, nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	fmt.Printf("Monitoring kubernetes result of all service (prometheus): %v\n", httpWriter.Body.String())

	httpRequest, err = http.NewRequest("GET", `http://172.17.8.100:30003/api/v1/query?query=kube_service_info%7Bservice%3D%22kubernetes%22%7D`, nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	fmt.Printf("Monitoring kubernetes result of kubernetes service (prometheus): %v\n", httpWriter.Body.String())

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/monitoring/services/kubernetes", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	fmt.Printf("Monitoring kubernetes result (API): %v\n", httpWriter.Body.String())

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/monitoring/services", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/monitoring/services?namespace=.*", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestGetServiceMetrics() {
	namespace := "vortex"
	services, err := suite.sp.KubeCtl.GetServices(namespace)
	suite.NoError(err)
	serviceName := services[0].GetName()

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/services/"+serviceName, nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestListControllerMetrics() {
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/controllers/", nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	httpRequest, err = http.NewRequest("GET", "http://localhost:7890/v1/monitoring/controllers?namespace=.*", nil)
	suite.NoError(err)

	httpWriter = httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}

func (suite *PrometheusTestSuite) TestGetControllerMetrics() {
	namespace := "vortex"
	deployments, err := suite.sp.KubeCtl.GetDeployments(namespace)
	suite.NoError(err)
	deploymentName := deployments[0].GetName()

	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/monitoring/controllers/"+deploymentName, nil)
	suite.NoError(err)

	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)
}
