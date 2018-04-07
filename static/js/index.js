
// polyfill fetch
require("whatwg-fetch");
const { createStoryServiceClient } = require("./chn_pb_twirp.js");
const useJSON = /json/.test(window.location.search);
const client = createStoryServiceClient("http://localhost:8080", useJSON);


function getStories(category) {

    client.stories({ category }).then(stories => {
        const links = document.getElementById('links');
        for (let story of stories.storiesList) {

            let entry = document.createElement('li');
            let title = document.createElement('a');
            title.classList.add('title');
            title.href = story.url;
            title.innerText = story.title;
            entry.appendChild(title);
            let score = document.createElement('score');
            score.classList.add('meta');
            score.innerText = story.score;
            entry.appendChild(score);
            links.appendChild(entry);
        }
     }, err => {

    });
}

getStories("best");