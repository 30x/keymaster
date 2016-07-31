package nginx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/30x/keymaster/nginx"
	"io/ioutil"
	"os"
)

var _ = Describe("TestConfig", func() {

	It("should return no error from good config", func() {
		tmpfile, err := writeConf(valid_conf)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(tmpfile.Name())
		err = nginx.TestConfig(tmpfile.Name())
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return an error from a bad config", func() {
		tmpfile, err := writeConf(invalid_conf)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(tmpfile.Name())
		err = nginx.TestConfig(tmpfile.Name())
		Expect(err).To(HaveOccurred())
	})

	It("should return an error when warnings", func() {
		tmpfile, err := writeConf(warn_conf)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(tmpfile.Name())
		err = nginx.TestConfig(tmpfile.Name())
		Expect(err).To(HaveOccurred())
	})
})

func writeConf(content string) (*os.File, error) {
	tmpfile, err := ioutil.TempFile("", "TestConfig")
	if err != nil {
		return nil, err
	}

	_, err = tmpfile.Write([]byte(content))
	if err != nil {
		return nil, err
	}
	err = tmpfile.Close()

	return tmpfile, err
}

var valid_conf = `
events {}
http {
  server {
      listen  9000;
      server_name  example.com;
      location / {
      }
  }
}`

var warn_conf = `
events {}
http {
  server {
      listen  9000;
      server_name  example.com;
      location / {
      }
  }
  server {
      listen  9000;
      server_name  example.com;
      location / {
      }
  }
}`

var invalid_conf = `
http {
  server {
      listen  9000;
      server_name  example.com;
      location / {
      }
  }
}`
