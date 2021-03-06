
function search() {
    if (searchText != '') {
        if (spells.length > 0) {
            spellSet = spells.filter(filterSearch);
            spellSet.sort(orderDoc);
        }
    } else {
        setSpellSet();
    }
    genResults(paginate(spellSet));
}

function filterSearch(doc) {
    var name = doc.data.Name.toLowerCase();
    return name.indexOf(searchText) >= 0;
}
