function binarySearch(array, target, lb, ub) {
	if (lb > ub) {
		return -1;
	}
	var mid = Math.floor((lb + ub) / 2);
	var val = array[mid];
	if (val > target) {
		return binarySearch(array, target, lb, mid - 1);
	}
	if (val < target) {
		return binarySearch(array, target, mid + 1, ub);
	}
	return mid;
}

function removeValue(array, target) {
	sorted = array.sort(sortNumber);
    var index = binarySearch(sorted, target, 0, array.length-1);
    var newList = sorted.slice(0, index);
    return newList.concat(sorted.slice(index + 1, sorted.length));
}

function sortNumber(a,b) {
    return a - b;
}

function setCookie(cname, cvalue, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + (exdays*24*60*60*1000));
    var expires = "expires="+d.toUTCString();
    document.cookie = cname + '=' + cvalue + '; ' + expires +'; path=/';
}

function setFlash(type, message) {
	var val = btoa(type + ':'+ message);
	while (val[val.length - 1] == '=') {
		val = val.slice(0, (val.length - 1));
	}
	setCookie('flash', val, 1);
}
