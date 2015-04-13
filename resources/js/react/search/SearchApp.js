import React from 'react';

var SearchTabControl = require('./SearchTabControl');
var SearchResultList = require('./SearchResultList');
var Loading = require('../Loading');
var Alert = require('../Alert');
var SearchLoadMore = require('./SearchLoadMore');
var $ = require('jquery');

class SearchApp extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            items: [],
            badges: {
                total: 0,
                users: 0,
                groups: 0,
                threads: 0
            },
            currentTab: 0,
            keyword: '',
            showResult: false,
            showError: false,
            error: '',
            page: 0,
            lastPage: true
        };
        this.cache = {
            keyword : ''
        }
    }

    handleSubmit(e) {
        e.preventDefault();
        if (this.state.keyword != "") {

            var $this = this;
            $.getJSON('/search/do?q='+this.state.keyword, function(response) {
                if (response.error) {
                    $this.setState({showError: true, showResult : false, error: response.error});
                    return;
                }
                var newState = {}
                newState.currentTab = 0;
                newState.items = response.results;
                newState.badges = {
                    total: response.total,
                    users: response.groups.users | 0,
                    groups: response.groups.groups | 0,
                    threads: response.groups.threads | 0
                };
                newState.page = response.page;
                newState.showError = false;
                newState.showResult = true;
                newState.error = '';

                if (Math.ceil(response.total / 10) <= response.page) {
                    newState.lastPage = true;
                } else {
                    newState.lastPage = false;
                }

                $this.cache.keyword = $this.state.keyword;
                $this.setState(newState);
            });
        }

    }

    changeHandler(e) {
        var keyword = e.target.value;
        this.setState({ keyword });
    }

    tabChangeHandler(e) {
        e.preventDefault();
        var selected = e.target.getAttribute("data-tab");
        if (selected != undefined && this.state.currentTab != selected) {
            var currentTab = parseInt(selected)
            this.setState({currentTab});
            var $this = this;
            var params = {q: this.cache.keyword, page: 1}
            if (selected > 0) {
                params["t"] = currentTab;
            }
            $.getJSON('/search/do', params, function(response) {
                if (response.error) {
                    $this.setState({showError: true, showResult : false, error: response.error});
                    return;
                }
                var newState = {}

                if (selected == 0 && response.page == 1) {
                    newState.badges = {
                        total: response.total,
                        users: response.groups.users | 0,
                        groups: response.groups.groups | 0,
                        threads: response.groups.threads | 0
                    };
                }
                newState.items = response.results;
                newState.page = response.page;
                newState.showError = false;
                newState.showResult = true;
                newState.error = '';
                if (Math.ceil(response.total / 10) <= response.page) {
                    newState.lastPage = true;
                } else {
                    newState.lastPage = false;
                }
                $this.setState(newState);
            });
        }
    }

    loadMoreHandler(e) {
        e.preventDefault();
        var $this = this;
        var params = {q: this.cache.keyword, page: this.state.page+1}
        if (this.state.currentTab > 0) {
            params["t"] = this.state.currentTab;
        }
        $.getJSON('/search/do', params, function(response) {
            if (response.error) {
                $this.setState({showError: true, showResult : false, error: response.error});
                return;
            }
            var newState = {}

            newState.items = $this.state.items.concat(response.results);
            newState.page = response.page;
            newState.showError = false;
            newState.showResult = true;
            newState.error = '';

            if (Math.ceil(response.total / 10) <= response.page) {
                newState.lastPage = true;
            } else {
                newState.lastPage = false;
            }
            $this.setState(newState);
        });
    }

    render() {
        return (
            <div className="search">
              <h1>搜索</h1>
              <p className="lead">
                  <form onSubmit={this.handleSubmit.bind(this)}>
                      <input type="text" className="form-control" placeholder="输入要查询的关键字" value={this.state.keyword} onChange={this.changeHandler.bind(this)} />
                  </form>
              </p>
              <Alert type="danger" className={! this.state.showError ? "hidden" : ""} text={this.state.error} />
              <SearchTabControl badges={this.state.badges} selected={this.state.currentTab} className={! this.state.showResult ? "hidden" : ""} onClick={this.tabChangeHandler.bind(this)}  />
              <SearchResultList items={this.state.items} className={! this.state.showResult ? "hidden" : ""} />
              <SearchLoadMore className={this.state.showResult && ! this.state.lastPage ? '' : 'hidden'} onClick={this.loadMoreHandler.bind(this)} />
              <Loading />
            </div>
        );
    }
}

export default SearchApp
