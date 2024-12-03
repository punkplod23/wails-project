import './style.css';
import './app.css';

import logo from './assets/images/logo-universal.png';
import {Greet} from '../wailsjs/go/main/App';
import {RunCSV} from '../wailsjs/go/main/App';
document.querySelector('#app').innerHTML = `
    <img id="logo" class="logo">
      <div class="result" id="result" style="min-height:40px; position:relative; width:600px; height:auto; display:block;">Please upload your file 👇</div>
      <div class="input-box" id="input" style='display:block; position:relative;'>
        <input class="input" style='display:none;' id="name" type="text" autocomplete="off" />
        <button class="btn" id="file">CSV</button>
        <button class="btn" id="search_button" onclick="greet()" style='display:none;'>Search</button>
      </div>
    </div>
`;
document.getElementById('logo').src = logo;

let nameElement = document.getElementById("name");


let fileElement = document.getElementById("file");
//nameElement.focus();
// Setup the greet function
window.greet = function () {
    // Get name
    
    let name = nameElement.value;
    console.log(name)
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

       
        nameElement.style.removeProperty('display');
        document.getElementById("file").style.display = 'none';
        document.getElementById("search_button").style.removeProperty('display'); 
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


