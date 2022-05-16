rm -rf dist
rm -rf hktdist
XX1=`pwd`
cd ~/MyWork/hackerToolsApp/app
yarn install
yarn build
mv dist $XX1/
git commit -m "up" .
git push

cd ~/MyWork/Hack-Tools
yarn install
yarn build
mv dist $XX1/hktdist
git commit -m "up" .
git push

cd $XX1
source ~/.zshrc
go get -u
go mod tidy
go build main.go
git commit -m "up" .
git push
