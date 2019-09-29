package cloud

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

//Client ...
type Client struct {
	URL      *url.URL
	Username string
	Password string
}

//Error struct
type Error struct {
	Exception string `xml:"exception"`
	Message   string `xml:"message"`
}

//Dial address parse
func Dial(host, username, password string) (*Client, error) {
	url, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	return &Client{
		URL:      url,
		Username: username,
		Password: password,
	}, nil
}

//Mkdir mkdir
func (c *Client) Mkdir(path string) error {
	_, err := c.sendRequest("MKCOL", path)
	return err
}

//sendRequest send request
func (c *Client) sendRequest(request, path string) ([]byte, error) {
	folderUrl, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest(request, c.URL.ResolveReference(folderUrl).String(), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Username, c.Password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(body) > 0 {
		error := Error{}
		err = xml.Unmarshal(body, &error)
		if err != nil {
			return body, fmt.Errorf("Error during XML Unmarshal for response %s. The error was %s", body, err)
		}
		if error.Exception != "" {
			return nil, fmt.Errorf("Exception: %s, Message: %s", error.Exception, error.Message)
		}
	}
	return body, nil
}

//Exists 检测是否存在
func (c *Client) Exists(path string) bool {
	_, err := c.sendRequest("PROPFIND", path)
	return err == nil
}

//Download download
func (c *Client) Download(path string) ([]byte, error) {
	pathUrl, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	//create the https request
	client := http.Client{}
	req, err := http.NewRequest("GET", c.URL.ResolveReference(pathUrl).String(), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.Password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	error := Error{}
	err = xml.Unmarshal(body, &error)
	if err == nil {
		if error.Exception != "" {
			return nil, fmt.Errorf("Exception: %s, Message: %s", error.Exception, error.Message)
		}
	}
	return body, nil
}

func (c *Client) Upload(src []byte, dest string) error {
	destUrl, err := url.Parse(dest)
	if err != nil {
		return err
	}
	//create the https request
	client := http.Client{}
	req, err := http.NewRequest("PUT", c.URL.ResolveReference(destUrl).String(), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(body) > 0 {
		error := Error{}
		err = xml.Unmarshal(body, &error)
		if err != nil {
			return fmt.Errorf("Error during XML Unmarshal for response %s. The error is %s", body, err)
		}
		if error.Exception != "" {
			return fmt.Errorf("Exception: %s, Message: %s", error.Exception, error.Message)
		}
	}
	return nil
}
