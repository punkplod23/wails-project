import './style.css';
import './app.css';

import logo from './assets/images/logo-universal.png';
import {Greet} from '../wailsjs/go/main/App';
import {RunCSV} from '../wailsjs/go/main/App';
document.querySelector('#app').innerHTML = `
    <img id="logo" class="logo">
      <div class="result" id="result">Please enter your name below 👇</div>
      <div class="input-box" id="input">
        <input class="input" id="name" type="text" autocomplete="off" />
        <button class="btn" id="file">CSV</button>
        <button class="btn" onclick="greet()">Greet</button>
      </div>
    </div>
`;
document.getElementById('logo').src = logo;

let nameElement = document.getElementById("name");


let fileElement = document.getElementById("file");
nameElement.focus();
// Setup the greet function
window.greet = function () {
    // Get name
    let name = nameElement.value;

    // Check if the input is empty
    if (name === "") return;

    // Call App.Greet(name)
    try {
        Greet(name)
            .then((result) => {
                // Update result with data back from App.Greet()
                resultElement.innerText = result;
            })
            .catch((err) => {
                console.error(err);
            });
    } catch (err) {
        console.error(err);
    }
};

let resultElement = document.getElementById("result");

fileElement.onclick = () => {

  try {
    window.go.main.App.SelectFile()
      .then((result) => {
        
        if (result === "") return;

        runCSVLocal(result);
      })
      .catch((err) => {
        console.error(err);
      });
  } catch (err) {
    console.error(err);
  }
}

function runCSVLocal(filePath){
    try {
        RunCSV(filePath)
            .then((result) => {

                // Update result with data back from App.Greet()
                resultElement.innerText = result;
            })
            .catch((err) => {
                console.error(err);
            });
    } catch (err) {
        console.error(err);
    }

}


