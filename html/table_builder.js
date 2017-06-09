
function buildHtmlTable(selector, rows, columns, changeFilterKeyColumn, changeFilterValueColumn) {
  currentRows = filterChanged(explodedItems(rows), changeFilterKeyColumn, changeFilterValueColumn)
  if (columns == null) {
	  columns = collectColumns(currentRows);
  }
  addColumns(selector, columns)
 
  for (var i = currentRows.length - 1; i >= 0; i--) {
    var row$ = $('<tr/>');
    for (var colIndex = 0; colIndex < columns.length; colIndex++) {
      var cellValue = currentRows[i][columns[colIndex]];
      if (cellValue == null) cellValue = "";
      row$.append($('<td/>').html(cellValue));
    }
    $(selector).append(row$);
  }
}

function addColumns(selector, columns) {
  var headerTr$ = $('<tr/>');
  for (var colIndex = 0; colIndex < columns.length; colIndex++) {
	headerTr$.append($('<th/>').html(columns[colIndex]));
  }
  $(selector).append(headerTr$);
}

function collectColumns(rows, selector) {
  var columnSet = [];
  
  for (var i = 0; i < rows.length; i++) {
    var rowHash = rows[i];
    for (var key in rowHash) {
      if ($.inArray(key, columnSet) == -1) {
		columnSet.push(key);	
      }
    }
  }
  return columnSet;
}


function explodedItems(rows) {
  var ret = []
  for (var i = 0; i < rows.length; i++) {
    var commonColumns = {}
    var rowHash = rows[i];
    var expanded = false
    for (var key in rowHash) {
		var cellValue = rowHash[key]
		if ($.isArray(cellValue)) {  
            expanded = true
			for (var ci = 0; ci < cellValue.length; ci++) {
				var childRow = cellValue[ci]
                var item = jQuery.extend({}, commonColumns);
                for (var childKey in childRow) {
					item[childKey] = childRow[childKey]
				}
                ret.push(item)
			}
		} else {
            commonColumns[key] = cellValue
        }
    }
    if (!expanded) {
        ret.push(rowHash)
    }
  }
  return ret
}

function filterChanged(rows, keyColumn, valueColumn) {
  var ret = []
  var lastStates = {}
  for (var i = 0; i < rows.length; i++) {
    var row = rows[i];
    var key = row[keyColumn]
    var value = row[valueColumn]
    if (lastStates[key] == null || lastStates[key] != value) {
        ret.push(row)
        lastStates[key] = value
    }
  }
  return ret
}