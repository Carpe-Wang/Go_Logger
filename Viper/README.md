#What is viper
Viper is a complete configuration solution for go applications. It is designed in applications and can handle all types of configuration requirements and formats.

# viper install
```go
go get github.com/spf13/viper
```

#why is viper 
>ADVANTAGES
```go
1:default setting
2:From Json TOML Yaml HCL Envfile and Java properties configuration file
3:Read and monitor configuration changes from the remote configuration system
```
# ViperFirstDemo.go

Simple Viper configuration and use process and startup.
```go
viper.SetDefault("filedir", "./")
	viper.SetConfigName("config")       //Configure file name and no extension
	viper.SetConfigType("yaml")         //Configure file type
	viper.AddConfigPath("/etc/appname") //Find configuration file path
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")
```

and config routine

```go
r := gin.Default()
	r.GET("/version", func(context *gin.Context) {
		context.String(http.StatusOK, viper.GetString("version"))
	})
	r.Run()
```
after RUNNING

We can see "v0.0.1" in the web localhost:8080/version. and we can know the "v0.0.1" is yaml file .
```
version : "v0.0.1"
```

in addtion, we can use viper to read byte's slice
```go
var yamlExample=[]byte('
hacker:true
name:carpe-wang
hobbies:
	women
        github
clothing
')
viper.ReadConfig(bytes.NewBuffer(yamlExample))
viper.Get("name")//可以等到carpe-wang
```