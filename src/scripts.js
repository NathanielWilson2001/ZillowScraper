const urlParams = new URLSearchParams(window.location.search);
const locationName = urlParams.get('location');
var styleElement = document.getElementById("image");
//styleElement.innerHTML = `<img src="`+locationName+`.jpg" alt="`+locationName+`" class="absolute inset-0 w-full h-full object-cover rounded-l-lg" loading="lazy" />`
       
var count;
var averagePrice;
var averageSquareFootSum;
var averagePricePerSquareFoot;
var averageRentZestimate;
var averageZestimate;
var condo;
var multiFamily;
var singleFamily;
 fetch('/fetchData', {
  method: 'POST',
  body: locationName
})
.then(function(response) {
  // Handle response from Go server
  if (response.ok) {
    console.log(response)
    return response.clone().json();
  } else {
    // Handle error response
    console.error('Error:', response.status);
  }
})
.then(data => {
  // Accessing specific elements from the JSON response
  /*
  console.log(data); // Using dot notation
  htmlToAdd = ``
  var styleElement = document.getElementById("Title");
  count = data.runningTotalEntries
  averagePrice = data.averagePriceSum
  averageSquareFootSum = data.averageSquareFootSum
  averagePricePerSquareFoot = data.averagePricePerSquareFoot
  averageRentZestimate = data.averageRentZestimate
  averageZestimate = data.averageZestimate
  condo = data.condo 
  multiFamily = data.multiFamily
  singleFamily  = data.singleFamily

 var styleElement = document.getElementById("Count");
 styleElement.innerHTML = count + " Houses Anaylized"

 var styleElement = document.getElementById("Title");
styleElement.innerHTML = "Typical House in " + locationName

var radios = document.getElementsByName["size"]

styleElement = document.getElementById("Data");
styleElement.textContent = "Is worth $" + averagePrice
*/

})
.catch(function(error) {
  // Handle network error
  console.error('Error:', error);
});

function changeOutput(value){
    styleElement = document.getElementById("Data");
    if(value == "averagePrice"){
    styleElement.textContent = "Is worth $" + averagePrice.toFixed(2)
    }
    if(value == "averageSquareFootage"){
        styleElement.textContent = "Is " + averageSquareFootSum.toFixed(2) + " square feet"
    }
    if(value == "Zillow"){
        styleElement.innerHTML = "Zillow's Rent Estimate is $" + averageRentZestimate.toFixed(2) + "<br>Zillow's Price Estimage is $" + averageZestimate.toFixed(2)
    }
}

   