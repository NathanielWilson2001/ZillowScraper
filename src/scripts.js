const urlParams = new URLSearchParams(window.location.search);
locationName = urlParams.get('location');   
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
  cityList = data.cityDataList
  switch(locationName){
    case "NorthEast":
        locationName = "North East"
    case "WestCoast":
        locationName = "West Coast"   
  }

  document.getElementById("body").innerHTML += 
  `<div class="bg-white mb-48 p-12 text-center text-neutral-700 top-0">
  <h2 class="mb-4 text-4xl font-semibold">`+locationName+` United States</h2>
  <h4 class="mb-6 text-xl font-semibold">Zillow Scraper</h4>
  <button
  type="button"
  data-te-ripple-init
  data-te-ripple-color="light"
  class="rounded bg-cyan-400 text-white  px-7 pb-2.5 pt-3 text-sm font-medium uppercase leading-normal  shadow-[0_4px_9px_-4px_#0891b2] transition duration-150 ease-in-out hover:bg-primary-600 hover:shadow-[0_8px_9px_-4px_#0891b2,0_4px_18px_0_#0891b2] focus:bg-primary-600 focus:shadow-[0_8px_9px_-4px_#0891b2,0_4px_18px_0_#0891b2] focus:outline-none focus:ring-0 active:bg-primary-700 active:shadow-[0_8px_9px_-4px_#0891b2,0_4px_18px_0_#0891b2] dark:shadow-[0_4px_9px_-4px_#0891b2] dark:hover:shadow-[0_8px_9px_-4px_#0891b2,0_4px_18px_0_#0891b2] dark:focus:shadow-[0_8px_9px_-4px_#0891b2,0_4px_18px_0_#0891b2] dark:active:shadow-[0_8px_9px_-4px_#0891b2,0_4px_18px_0_#0891b2]"  onclick="returnHome()">
  Back to Home
</button>
  <button
    type="button"
    data-te-ripple-init
    data-te-ripple-color="light"
    class="disc
    ">
    Regional Trends
  </button>
  <div id="trends"></div>
</div>`;
  var count = 0;
  data.cityDataList.forEach(element => {
   console.log(element.data);
   document.getElementById("body").innerHTML += 
   `<div class="flex justify-center items-center my-10">
    <div class="flex font-sans shadow-2xl w-7/12 h-10/12">
      <div class="flex-none w-1/3 relative" id="image"><img src="images/`+element.location+`.jpg" alt="`+element.location+`" class="absolute inset-0 w-full h-full object-cover rounded-l-lg" loading="lazy" /></div>
          <form class="flex-auto p-6  bg-white" name="dataValues`+count+`">
                <div class="flex flex-wrap">
                  <h1 class="flex-auto font-medium text-2xl text-slate-900" id="Count">`+ element.data.runningTotalEntries + ` Houses Anaylized!</h1>
                  <div class="w-full flex-none mt-2 order-1 text-3xl font-bold text-cyan-600" id="Title">Typical House in `+element.location+`</div>
                </div>
                <div class="flex items-baseline mt-4 mb-6 pb-6 border-b border-slate-200">
                  <div class="space-x-2 flex text-sm font-bold">
                    <label>
                      <input class="sr-only peer" name="size" type="radio" onclick="changeOutput('averagePrice', '`+count+`')" checked/>
                      <div class="w-16 h-16 rounded-full flex items-center justify-center text-cyan-400 text-2xl peer-checked:text-white peer-checked:bg-cyan-400 peer-checked:shadow-[0_8px_9px_-4px_#0891b2,0_4px_18px_0_#0891b2] peer-checked:outline-none peer-checked:ring-0><span class="text-2xl">&#36;</span></div>
                    </label>
                    <label>
                      <input class="sr-only peer" name="size" type="radio" value="s"  onclick="changeOutput('averageSquareFootage', '`+count+`')"/>
                      <div class="w-16 h-16 rounded-full flex items-center justify-center text-cyan-400 text-2xl peer-checked:text-white peer-checked:bg-cyan-400 peer-checked:shadow-[0_8px_9px_-4px_#0891b2,0_4px_18px_0_#0891b2] peer-checked:outline-none peer-checked:ring-0><span class="text-2xl">ft<sup>2</sup></span></div>
                    </label>
                    <label>
                      <input class="sr-only peer" name="size" type="radio"  onclick="changeOutput('Zillow', '`+count+`')"/>
                      <div class="w-16 h-16 rounded-full flex items-center justify-center text-cyan-400 text-xl peer-checked:text-white peer-checked:bg-cyan-400 peer-checked:shadow-[0_8px_9px_-4px_#0891b2,0_4px_18px_0_#0891b2] peer-checked:outline-none peer-checked:ring-0><span class="text-lg font-bold">Zillow</span></div>
                    </label>
                    <label>
                      <input class="sr-only peer" name="size" type="radio"  onclick="changeOutput('Trends', '`+count+`')"/>
                      <div class="w-16 h-16 rounded-full flex items-center justify-center text-cyan-400 text-xl peer-checked:text-white peer-checked:bg-cyan-400 peer-checked:shadow-[0_8px_9px_-4px_#0891b2,0_4px_18px_0_#0891b2] peer-checked:outline-none peer-checked:ring-0><span class="text-lg font-bold">Type</span></div>
                    </label>
                    </div>
                  </div>
                  <div class="w-max flex space-x-4 mb-5 text-sm font-medium h-1/2"> <div class="w-max flex-auto flex space-x-4 mr-4 mt-4">
                  <p class="w-max text-black tracking-wide text-xl antialiased proportional-nums" id="Data`+count+`">Is worth $`+cityList[count].data.averagePriceSum.toFixed(2)+`</p>
                </div>
              </div>
            </form>
          </div>
      </div>`;
      count++; 
  });
})
.catch(function(error) {
  // Handle network error
  console.error('Error:', error);
});

function changeOutput(value, count){
    styleElement = document.getElementById("Data"+count);
    if(value == "averagePrice"){
    styleElement.innerHTML = `<h4 class="text-xl">Is worth</h4> <p>$` + cityList[count].data.averagePriceSum.toFixed(2)+`</p>`
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
          text: "Housing Type Makeup",
          fontSize: 24,
          padding:5
        },
      }
    });
    }
}

function returnHome() {
    // Send an HTTP GET request to the Go backend if needed
    window.location.href = '/'
};
