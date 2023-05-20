document.addEventListener("DOMContentLoaded", function() {
  var selectedLocation = "default";
  document.getElementById("buttonBoston").addEventListener("click", function() {
    selectedLocation = "boston";
    // Send an HTTP GET request to the Go backend if needed
    fetch("/hello");

    // Redirect to the appropriate location page
    window.location.replace("localhost:5000/city.html?location=" + selectedLocation);
  });

  const urlParams = new URLSearchParams(window.location.search);
  const locationName = urlParams.get('location');
  if (locationName === "Boston") {
    alert("gere");
    var styleElement = document.getElementById("BackgroundImage");
    styleElement.innerHTML = `
      .background-image {
        background-image: url('boston.jpg');
        background-position: center;
        background-repeat: no-repeat;
        background-size: cover;
      }
    `;
  }
});
