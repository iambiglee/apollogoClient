package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/apollogoClient/v1/env"
	"github.com/apollogoClient/v1/env/config"
	"github.com/apollogoClient/v1/extension"
	"io/ioutil"
	"net"
	"net/http"
	url2 "net/url"
	"strings"
	"sync"
	"time"
)

var (
	//for on error retry
	onErrorRetryInterval = 2 * time.Second //2s

	connectTimeout = 1 * time.Second //1s

	//max retries connect apollo
	maxRetries = 5

	//defaultMaxConnsPerHost defines the maximum number of concurrent connections
	defaultMaxConnsPerHost = 512
	//defaultTimeoutBySecond defines the default timeout for http connections
	defaultTimeoutBySecond = 1 * time.Second
	//defaultKeepAliveSecond defines the connection time
	defaultKeepAliveSecond = 60 * time.Second
	// once for single http.Transport
	once sync.Once
	// defaultTransport http.Transport
	defaultTransport *http.Transport
)

// CallBack 请求回调函数
// TODO 这里面又是接口的又是回调的，是什么东西？第三行是个什么意思?
type CallBack struct {
	SuccessCallBack   func([]byte, CallBack) (interface{}, error)
	NotModifyCallBack func() error
	AppConfigFunc     func() config.AppConfig
	Namespace         string
}

// Request 建立网络请求,这里只是建立的IP连接，长连接？
// TODO 第二行的等于接口指针是什么？url2是什么思想感情
//TODO strings 和 url 包是什么
//TODO callBack 是什么
func Request(requestURL string, connectionConfig *env.ConnectConfig, callback *CallBack) (interface{}, error) {
	client := &http.Client{}
	if connectionConfig != nil && connectionConfig.Timeout != 0 {
		client.Timeout = connectionConfig.Timeout
	} else {
		client.Timeout = connectTimeout
	}
	var err error
	url, err := url2.Parse(requestURL)
	if err != nil {
		fmt.Errorf("request Apollo Server error:%s", requestURL)
		return nil, err
	}
	var insecureSkipVerify bool
	if strings.HasPrefix(url.Scheme, "https") {
		insecureSkipVerify = true
	}
	client.Transport = getDefaultTransport(insecureSkipVerify)
	retry := 0
	var retries = maxRetries
	if connectionConfig != nil && !connectionConfig.IsRetry {
		retry = 1
	}
	for {
		retry++
		if retry > retries {
			break
		}
		req, err := http.NewRequest("GET", requestURL, nil)
		if req == nil || err != nil {
			fmt.Errorf("generate connect Apollo request Fail,url:%s,Error:%s", requestURL, err)
			return nil, errors.New("generate connect Apollo request fail")
		}

		//添加head 选项
		httpAuth := extension.GetHttPAuth()
		if httpAuth != nil {
			headers := httpAuth.HTTPHeaders(requestURL, connectionConfig.AppID, connectionConfig.Secrct)
			if len(headers) > 0 {
				req.Header = headers
			}
			host := req.Header.Get("Host")
			if len(host) > 0 {
				req.Host = host
			}
		}
		res, err := client.Do(req)
		if res != nil {
			defer res.Body.Close()
		}

		if res == nil || err != nil {
			// if error then sleep
			time.Sleep(onErrorRetryInterval)
			continue
		}

		//处理这种返回结果
		//TODO 这个下面的callback 是怎么回事
		switch res.StatusCode {
		case http.StatusOK:
			responseBody, err := ioutil.ReadAll(req.Body)
			if err != nil {
				// if error then sleep
				time.Sleep(onErrorRetryInterval)
				continue
			}
			if callback != nil && callback.SuccessCallBack != nil {
				return callback.SuccessCallBack(responseBody, *callback)
			}
			return nil, nil
		case http.StatusNotModified:
			if callback != nil && callback.NotModifyCallBack != nil {
				return nil, callback.NotModifyCallBack()
			}
		default:
			time.Sleep(onErrorRetryInterval)
			continue
		}

	}
	if retry > retries {
		err = errors.New("Over Max retry still Error")
	}
	return nil, err
}

//getDefaultTransport 这里又是做什么的，获取http 连接？
//TODO DialContext 的写法怎么这么夸张
func getDefaultTransport(insecureSkipVerify bool) *http.Transport {
	once.Do(func() {
		defaultTransport = &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			MaxIdleConns:        defaultMaxConnsPerHost,
			MaxIdleConnsPerHost: defaultMaxConnsPerHost,
			DialContext: (&net.Dialer{
				KeepAlive: defaultKeepAliveSecond,
				Timeout:   defaultTimeoutBySecond,
			}).DialContext,
		}
		if insecureSkipVerify {
			defaultTransport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: insecureSkipVerify,
			}
		}
	})
	return defaultTransport
}
