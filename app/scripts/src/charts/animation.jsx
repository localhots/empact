var ChartAnimationMixin = {
    animDuration: 350,

    animate: function(ref, attr, from, to) {
        var node = ref.getDOMNode(),
            anim = anim = document.createElementNS('http://www.w3.org/2000/svg', 'animate');

        anim.setAttributeNS(null, 'attributeType', 'XML');
        anim.setAttributeNS(null, 'attributeName', attr);
        anim.setAttributeNS(null, 'values', from +';'+ to);
        anim.setAttributeNS(null, 'dur', this.animDuration +'ms');
        anim.setAttributeNS(null, 'calcMode', 'spline');
        anim.setAttributeNS(null, 'keySplines', this.easing);
        anim.setAttributeNS(null, 'repeatCount', '1');
        node.appendChild(anim);
        anim.beginElement();
        setTimeout(function() {
            node.setAttributeNS(null, attr, to);
        }, this.animDuration);
    },

    clearAnimations: function(ref) {
        var node = ref.getDOMNode();
        while (node.firstChild) {
            node.removeChild(node.firstChild);
        }
    }
};
