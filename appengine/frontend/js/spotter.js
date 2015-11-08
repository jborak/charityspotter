

function searchItems(event) {

  event.preventDefault();

  var input = $("#search-bar")[0].value;
  if (typeof input == "undefined" || input === null) {
    return;
  }
  if (input.length == 0) {
    return;
  }

  var xhr =  new XMLHttpRequest();
  xhr.onreadystatechange = function() {
    if (4 != xhr.readyState) {
      return;
    } else if (200 != xhr.status) {
      return;
    }

    try {
      var response = JSON.parse(xhr.responseText)
      if (typeof response == "undefined" || response === null) {
        return;
      }

      console.log(response);
      displayResults(input, response);
    } catch (e) {
      console.log(e)
    }
  }

  xhr.open("POST", "http://charityspotter.com/api/search/", true);
  xhr.setRequestHeader("Content-type", "application/json");
  xhr.setRequestHeader("Access-Control-Allow-Origin", "*");

  var data = {"terms": input};
  xhr.send(JSON.stringify(data));
}

function displayResults(terms, items) {
  var searchResults = $("#results")[0];
  $(searchResults).empty();

  for (var i = 0; i < items.length; i++) {
    var title = ""; // items[i].data;
    var image = items[i].url; //"img/clothes-images/pants.png";
    var description = items[i].data; // "cool item #" + i;

    var item = document.createElement('div');
    item.className = 'col-md-3 result-box';
    item.innerHTML =   
      '<div class="intra-box">' +
      '    <h4>' + title + '</h4>' +
      '    <img class="image-preview" src="' + image + '">' +
      '    <p>' + description + '</p>' +
      '</div>';
    $(searchResults).append(item);
    // '<div class="col-md-3 result-box">' +
    // '</div>';
  }

  var searchTerms = $("#search-summary")[0];
  searchTerms.innerHTML = "Results for: &quot;" + terms + "&quot;";
  searchTerms.style.display = "block";

  var displayCount = $("#result-count")[0];
  displayCount.innerHTML = "Displaying " + items.length + " out of " + items.length + " items found";
  displayCount.style.display = "block";

  displayCount.parentNode.style.display = "block";
}
