var ChartAnimationMixin = {
    animDuration: 350,

    animate: function(ref, attr, from, to) {
        var node = ref.getDOMNode(),
            anim = document.createElementNS('http://www.w3.org/2000/svg', 'animate');

        anim.setAttributeNS(null, 'attributeType', 'XML');
        anim.setAttributeNS(null, 'attributeName', attr);
        anim.setAttributeNS(null, 'values', from +';'+ to);
        anim.setAttributeNS(null, 'dur', this.animDuration +'ms');
        anim.setAttributeNS(null, 'calcMode', 'spline');
        // Easings to consider:
        // easeOutCirc: 0.075 0.82 0.165 1
        // easeOutBack: 0.175 0.885 0.32 1.275
        // easeInOutCubic: 0.645 0.045 0.355 1
        anim.setAttributeNS(null, 'keySplines', '0.175 0.885 0.32 1.275');
        anim.setAttributeNS(null, 'repeatCount', '1');
        anim.addEventListener('endEvent', function() {
            node.setAttributeNS(null, attr, to);
        });
        node.appendChild(anim);

        anim.beginElement();
    },

    clearAnimations: function(ref) {
        var node = ref.getDOMNode();
        while (node.firstChild) {
            node.removeChild(node.firstChild);
        }
    }
};
