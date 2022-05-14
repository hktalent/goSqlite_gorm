rm -rf dist
XX1=`pwd`
cd ~/MyWork/hackerToolsApp/app
yarn build
mv dist $XX1/
cd $XX1
go get -u
go mod tidy
go build

