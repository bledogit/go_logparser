FILES=$(wildcard ./*.go) $(wildcard cmd/$(logparser)/*.go)
PACKAGE=logparser.zip
REGION=us-east-1
TARGET=logparser_lambda
ROLE=arn:aws:iam::254709965617:role/cloudapi1dot0/ADFS-WatchProd-CLOUDAPI1DOT0-LambdaExecutionRole-1RTLSRWSYAYPI
FUNCTION=cloudapi1dot0-logparser-redirect
TAGS=ApplicationName="CLOUD-API-1.0",CostCenter="5501983"
ENV=Variables="{ENDPOINT=http://sandbox-cloudapi.imrworldwide.com/,MAX_REQUESTS=1000,MAX_WORKERS=500}"

test:
	@echo Running test
	go test ./...

$(TARGET): $(FILES)
	@echo Building ... $@
	GOOS=linux go build -o $(TARGET) cmd/$(TARGET)/main.go

package: logparser_lambda
	@echo Building deployment package
	zip $(PACKAGE) logparser_lambda

create-function: 
	aws lambda  create-function \
	--region $(REGION) \
	--function-name $(FUNCTION) \
	--memory-size 128 \
	--role $(ROLE) \
	--runtime go1.x \
	--tags $(TAGS) \
	--zip-file  fileb://$(PACKAGE) \
	--handler $(TARGET)

redeploy: package
	aws lambda update-function-code \
	--function-name $(FUNCTION) \
	--zip-file  fileb://$(PACKAGE)

update: 
	aws lambda  update-function-configuration \
	--region $(REGION) \
	--function-name $(FUNCTION) \
	--memory-size 256 \
	--role $(ROLE) \
	--runtime go1.x \
	--environment $(ENV) \
	--handler $(TARGET)


	
clean:
	go clean
	rm -f $(PACKAGE)

all: package
default: all