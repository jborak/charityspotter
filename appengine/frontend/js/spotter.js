

function searchItems() {
  var input = $("#searchField")[0].value;
  if (typeof input == "undefined" || input === null) {
    return;
  }

  var xhr =  new XMLHttpRequest();
  http.onreadystatechange = function() {
    if (4 != http.readyState) {
      return;
    } else if (200 != http.status) {
      return;
    }

    try {
      var response = JSON.parse(responseText)
      if (typeof response == "undefined" || response === null) {
        return;
      }

      console.log(response);
      displayResults(response);
    } catch (e) {
      console.log(e)
    }
  }

  xhr.open("POST", "http://charityspotter.com/api/search/", true);
  xhr.setRequestHeader("Content-type", "application/json");
  xhr.send({"terms": value});
}

function displayResuls(items) {
  var searchResults = $("#searchResults")[0];
  $(searchResults).
  
  for (var i = 0; i < items.length; i++) {

  }
}

