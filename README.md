# Project Spore

A ```toy``` based on [baoz.cn](http://baoz.cn)

![screenshot](https://cloud.githubusercontent.com/assets/1967804/7112544/3e2953b6-e1fd-11e4-9e11-a4928ee333f3.png)


Goals:

- [x] Storage
    * [x] Mysql Migration
    * [x] Redis
- [x] Commands
    * [x] Crawler
    * [x] Server
    * [x] Stat
- [x] Search
    * [x] Coreseek integrated
- [x] Multi-dimensional Ranking
- [x] Web App
    * [x] React

> Project won't include features that might jeopardize personal infomation, such as Login or Register right now, unless Baoz support Official API with OAuth.

:smiley:

## How to play

- Software
    These are MUST HAVE.
    - Go
    - CoreSeek4+
    - NodeJS(with NPM)
    - Mysql
    - Redis

- Install

  ```
  go get github.com/chekun/spore
  cd $GOPATH/src/github.com/chekun/spore
  go get ...
  npm install
  ```
- Configure

  ```
  cp config.example.conf config.conf
  cp sphinx.example.conf sphinx.conf
  ```

  And edit to your own configs.

- Database Migration

  ```
  # create your database
  sql-migrate -env=development -config=config.yml
  ```

- Build And Run

  ```
  cd $GOPATH/src/github.com/chekun/spore
  gulp
  cd spored
  go build
  ```
- Crawl Data

  ```
  cd $GOPATH/src/github.com/chekun/spore/spored
  ./spored crawl -env=development -config=../config.yml
  ```

  This is process is totally automatic, you don't need to worry about the stop logic, we have that covered, Just put it in daemon mode.

- Run Stats

  Once data has fetched. You can run stats command

  ```
  ./spored stat -env=development -config=../config.yml
  ```

- Web App

  ```
  ./spored serve -env=development -config=../config.yml
  ```

  Open your browser, visit http://127.0.0.1:10999/ and have fun.


# The MIT License (MIT)
  Copyright (c) 2012 [@chekun](https://github.com/chekun)

  Permission is hereby granted, free of charge, to any person obtaining a copy
  of this software and associated documentation files (the “Software”), to deal
  in the Software without restriction, including without limitation the rights
  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
  copies of the Software, and to permit persons to whom the Software is
  furnished to do so, subject to the following conditions:

  The above copyright notice and this permission notice shall be included in
  all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
  THE SOFTWARE.
