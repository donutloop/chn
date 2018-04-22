// polyfill fetch
require("whatwg-fetch");
const { createStoryServiceClient } = require("./chn_pb_twirp.js");
const useJSON = /json/.test(window.location.search);
const client = createStoryServiceClient("http://localhost:8080", useJSON);

var app = new Vue({
    el: '#app',
    data: {
        stories: [],
        loading: false,
    },
    created: function() {
        this.loading = true;
        client.stories({ category:"best", limit: 50 }).then(stories => {
            this.stories = stories.storiesList
            this.loading = false;
        }, err => {
            this.loading = false;
        });
    },
})

