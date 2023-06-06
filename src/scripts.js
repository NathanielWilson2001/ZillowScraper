const urlParams = new URLSearchParams(window.location.search);
const locationName = urlParams.get('location');   
var cityList;
 fetch('/fetchData', {
  method: 'POST',
  body: locationName
})
.then(function(response) {
  // Handle response from Go server
  if (response.ok) {
    console.log(response.body)
    return response.clone().json();
  } else {
    // Handle error response
    console.error('Error:', response.status);
  }
})
.then(data => {
  console.log(data.cityDataList[0].location)
  cityList = data.cityDataList;
  var htmlToOutput = ""; 
  var count = 0;
  data.cityDataList.forEach(element => {
   console.log(element.data);
   htmlToOutput = htmlToOutput + 
   `<div class="flex justify-center items-center my-10">
    <div class="flex font-sans shadow-2xl w-7/12 h-9/12">
      <div class="flex-none w-96 relative" id="image"><img src="images/`+element.location+`.jpg" alt="`+element.location+`" class="absolute inset-0 w-full h-full object-cover rounded-l-lg" loading="lazy" /></div>
          <form class="flex-auto p-6  bg-white" name="dataValues`+count+`">
            <div class="flex flex-wrap">
              <h1 class="flex-auto font-medium text-2xl text-slate-900" id="Count">`+ element.data.runningTotalEntries + ` Houses Anaylized!</h1>
              <div class="w-full flex-none mt-2 order-1 text-3xl font-bold text-cyan-600" id="Title">Typical House in `+element.location+`</div>
            </div>
            <div class="flex items-baseline mt-4 mb-6 pb-6 border-b border-slate-200">
              <div class="space-x-2 flex text-sm font-bold">
                <label>
                  <input class="sr-only peer" name="size" type="radio" onclick="changeOutput('averagePrice', '`+count+`')" checked/>
                  <div class="w-14 h-14 rounded-full flex items-center justify-center text-cyan-400 peer-checked:bg-cyan-600 peer-checked:text-white"><span class="text-2xl">&#36;</span></div>
                </label>
                <label>
                  <input class="sr-only peer" name="size" type="radio" value="s"  onclick="changeOutput('averageSquareFootage', '`+count+`')"/>
                  <div class="w-14 h-14 rounded-full flex items-center justify-center text-cyan-400 peer-checked:bg-cyan-600 peer-checked:text-white"><span class="text-2xl">ft<sup>2</sup></span></div>
                </label>
                <label>
                  <input class="sr-only peer" name="size" type="radio"  onclick="changeOutput('Zillow', '`+count+`')"/>
                  <div class="w-14 h-14 rounded-full flex items-center justify-center text-cyan-400 peer-checked:bg-cyan-600 peer-checked:text-white"><span class="text-lg font-bold">Zillow</span></div>
                </label>
                 <label>
                  <input class="sr-only peer" name="size" type="radio"  onclick="changeOutput('Trends', '`+count+`')"/>
                  <div class="w-14 h-14 rounded-full flex items-center justify-center text-cyan-400 peer-checked:bg-cyan-600 peer-checked:text-white"><span class="text-lg font-bold">Trends</span></div>
                </label>
                </div>
              </div>
              <div class="flex space-x-4 mb-5 text-sm font-medium h-1/2"> <div class="flex-auto flex space-x-4 mt-4">
              <p class="text-black tracking-wide text-xl antialiased proportional-nums" id="Data`+count+`">Is worth $`+cityList[count].data.averagePriceSum.toFixed(2)+`</p>
            </div>
          </div>
            </form>
          </div>
      </div>`;
      count++; 
  });
  document.getElementById("body").innerHTML = htmlToOutput;
})
.catch(function(error) {
  // Handle network error
  console.error('Error:', error);
});

function changeOutput(value, count){
    styleElement = document.getElementById("Data"+count);
    if(value == "averagePrice"){
    styleElement.innerHTML = "Is worth $" + cityList[count].data.averagePriceSum.toFixed(2)
    }
    if(value == "averageSquareFootage"){
        styleElement.textContent = "Is " + cityList[count].data.averageSquareFootSum.toFixed(2) + " square feet"
    }
    if(value == "Zillow"){
        styleElement.innerHTML = "Zillow's Rent Estimate is $" + cityList[count].data.averageRentZestimate.toFixed(2) + "<br>Zillow's Price Estimage is $" + cityList[count].data.averageZestimate.toFixed(2)
    }
    if(value == "Trends"){
      styleElement.innerHTML = `<canvas class="text-xl" id="chart`+count+`" style="max-width:600px;height:400px;font-size: x-large;"></canvas>`
      var xValues = ["Single Family", "Multi-Family", "Condo"];
    var yValues = [cityList[count].data.singleFamily, cityList[count].data.multiFamily, cityList[count].data.condo];
    var barColors = [
      "#083344",
      "#0891b2",
      "#22d3ee"
    ];
    
    var chartName = "chart"+count
    new Chart(chartName, {
      type: "pie",
      data: {
        labels: xValues,
        datasets: [{
          backgroundColor: barColors,
          data: yValues
        }]
      },
      options: {
        title: {
          display: true,
          text: "Housing type makeup",
          fontSize: 24,
          padding:5
        },
      }
    });

    }
}
