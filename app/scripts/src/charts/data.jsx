var ChartDataMixin = {
    apiParams: function() {
        var params = _.clone(this.props.params);
        params['item'] = this.state.item;
        return params;
    },

    apiParamsHash: function() {
        var params = this.apiParams();
        var pairs = _.chain(params)
            .keys()
            .sort()
            .map(function(param) {
                return [param, params[param]];
            })
            .value();
        return JSON.stringify(pairs);
    },

    fetchData: function() {
        var paramsHash = this.apiParamsHash();

        if (!this.state.item) {
            return;
        }
        if (this.state.currentApi === this.props.api && this.state.currentParams === paramsHash) {
            return;
        }

        // console.log('-----> fetching', this.props.api, this.state.item);
        this.setState({
            currentApi: this.props.api,
            currentParams: paramsHash,
            state: 'loadingData'
        }, function() {
            $.get(this.props.api, this.apiParams(), function(res){
                this.setState({
                    rawData: res,
                    state: 'newData'
                }, this.handleNewData);
            }.bind(this));
        }.bind(this));
    }
};
