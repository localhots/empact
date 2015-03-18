var InfoBlock = React.createClass({
    render: function() {
        return (
            <div className="info-block">
                <div className={'img'+ (this.props.image ? '' : ' empty') +' '+ this.props.className}
                    style={this.props.image ? {backgroundImage: "url("+ this.props.image +")"} : null} />
                <h1>{this.props.title}</h1>
            </div>
        )
    }
});
