var pageSize = 10;
var page = 0;
var pages;
var test;

function paginate(dataSet) {
    pages = Math.ceil(dataSet.length / pageSize);
    var lb;
    var ub;
    ub = ((pages - page) >= 5) ? page + 4 : pages - 1 ;

    if (page < 6) {
		ub = (pages > 9) ? 9 : pages - 1
	}
    lb = (((ub - 9) > 0) ? ub - 9 : 0);
    generatePaginator(lb, ub, pages);
    var beg = page * pageSize;
    var end = ((page * pageSize) + pageSize);
    return dataSet.slice(beg, end);
}

function generatePaginator(lb, ub, pages) {
    var paginator = $('#paginator');
    paginator.html('');
    if (pages < 2) {
        return
    }
    var prev = $('<li id="prev"><a aria-label="Previous" style="cursor:pointer;"><span aria-hidden="true">&laquo;</span</a></li>');
    var next = $('<li id="next"><a aria-label="Next" style="cursor:pointer;"><span aria-hidden="true">&raquo;</span></a></li>');
    if (page == 0) {
        prev.addClass('disabled');
    }
    if (page === (pages - 1)) {
        next.addClass('disabled');
    }
    prev.click(function() {
        if (page > 0) {
            page--;
            genResults(paginate(spellSet));
        }
    });
    next.click(function() {
        if (page + 1 < pages) {
            page++;
            genResults(paginate(spellSet));
        }
    });
    paginator.append(prev);
    for (var i =  lb; i <= ub; i++) {
        var elem = $('<li data-page="' + i + '"><a href="#">' + (i+1) + '</a></li>');
        elem.addClass('hidden-xs');
        if (page == i) {
            elem.addClass('active');
            elem.removeClass('hidden-xs');
        }
        elem.click(function() {
            page =+ this.getAttribute('data-page');
            genResults(paginate(spellSet));
        });
        paginator.append(elem);
    }
    paginator.append(next);
}

function genResults(filteredSpells) {
    var i;
    group = $('<div class="panel-group" id="dndSpellAccordion" role="tablist"></div>');
    for (i = 0; i < filteredSpells.length; i++) {
        var doc = filteredSpells[i]
        spell = $('<div class="panel panel-default">'+
                '<div class="panel-heading clearfix" role="tab" id="heading' + doc.id + '">'+
                    '<h4 class="panel-title">'+
                        '<a data-toggle="collapse" data-parent="#dndSpellAccordion" data-target="#dndSpell' + doc.id + '" style="cursor:pointer;">'+
                            doc.data.Name +
                        '</a>'+
                        '<span class="visible-xs-block"><br></span>'+
                        '<form class="addSpell col-lg-6 pull-right">' +
                            ((doc.data.Custom && cat == 'userC') ? '<div class="col-lg-2 pull-right"><a href="/edit/spell/' + doc.id + '" class="btn btn-sm btn-primary">Edit</a></div>' : '') +
                            '<div class="input-group pull-right input-group-sm col-lg-10 col-xs-8">'+
                                '<input type="number" name="spellLvl" class="form-control" placeholder="Level" required>'+
                                '<span class="input-group-btn">'+
                                    '<button class="btn btn-default">Add Spell</button>'+
                                '</span>'+
                            '</div>'+
                            '<input type="hidden" name="spellId" value="' + doc.id + '">' +
                            '<input type="hidden" name="userId" value="' + userId + '"/>' +
                        '</form>' +
                    '</h4>'+
                '</div>'+
                '<div id="dndSpell' + doc.id + '" class="panel-collapse collapse">'+
                    '<div class="panel-body description-panel">' +
                        '<h4>' + doc.data.School +
                            ((doc.data.Subschool != '') ? '<span> (' + doc.data.Subschool + ')</span>' : '') +
                            ((doc.data.Descriptors != '') ? '<span> [' + doc.data.Descriptors + ']</span>' : '') +
                        '</h4>' +
                        ((doc.data.Rulebook != '') ? '<p>' +
                            doc.data.Rulebook + ((doc.data.Page != 0 ) ? '<span> p. ' + doc.data.Page + '</span>' : '') +
                        '</p>' : '') +
                        '<table>' +
                            '<tbody>' +
                                ((doc.data.Components != '') ? '<tr><td class="text-right"><strong>Components:</strong> &nbsp;</td><td class="text-left">' + doc.data.Components + '</td></tr>' : '') +
                                ((doc.data.Displays != '') ? '<tr><td class="text-right"><strong>Displays:</strong> &nbsp;</td><td class="text-left">' + doc.data.Displays + '</td></tr>' : '') +
                                ((doc.data.CastingTime != '') ? '<tr><td class="text-right"><strong>Casting Time:</strong> &nbsp;</td><td class="text-left">' + doc.data.CastingTime + '</td></tr>' : '') +
                                ((doc.data.SpellRange != '') ? '<tr><td class="text-right"><strong>Range:</strong> &nbsp;</td><td class="text-left">' + doc.data.SpellRange + '</td></tr>' : '') +
                                ((doc.data.Area != '') ? '<tr><td class="text-right"><strong>Area:</strong> &nbsp;</td><td class="text-left">' + doc.data.Area + '</td></tr>' : '') +
                                ((doc.data.Effect != '') ? '<tr><td class="text-right"><strong>Effect:</strong> &nbsp;</td><td class="text-left">' + doc.data.Effect + '</td></tr>' : '') +
                                ((doc.data.Target != '') ? '<tr><td class="text-right"><strong>Target:</strong> &nbsp;</td><td class="text-left">' + doc.data.Target + '</td></tr>' : '') +
                                ((doc.data.Duration != '') ? '<tr><td class="text-right"><strong>Duration:</strong> &nbsp;</td><td class="text-left">' + doc.data.Duration + '</td></tr>' : '') +
                                ((doc.data.SavingThrow != '') ? '<tr><td class="text-right"><strong>Saving Throw:</strong> &nbsp;</td><td class="text-left">' + doc.data.SavingThrow + '</td></tr>' : '') +
                                ((doc.data.SpellResistance != '') ? '<tr><td class="text-right"><strong>Spell Resistance:</strong> &nbsp;</td><td class="text-left">' + doc.data.SpellResistance + '</td></tr>' : '') +
                            '</tbody>' +
                        '</table>' +
                        '<br>' +
                        '<span>' +
                            '<div ' + ((doc.data.Custom) ? 'class="description"' : '') + '>' + doc.data.DescriptionHtml + '</div>' +
                        '</span>' +
                    '</div>' +
                '</div>'+
            '</div>');
        group.append(spell);
    }
    $('#spells').html(group);
    $('form.addSpell').on("submit", function(e) {
        e.preventDefault();
        var lvlStr = $(this).find('input[name="spellLvl"]').val();
        var lvl =+ $(this).find('input[name="spellLvl"]').val();
        if (lvl < 0 || lvl > 9 || lvlStr === '') {
            console.log("ERROR");
            return
        }
        $.ajax({
            type: 'POST',
            url: '/user/addSpell',
            data: $(this).serialize(),
            success: function(data) {
                try {
                    data = JSON.parse(data);
                    if (data.success) {
                        picked = data.picked;
                        renderPicked();
                    }
                } catch(err) {
                    setFlash('alertError', 'Your session has expired. Please login');
                    window.location.pathname = '/'
                }
            },
        });

        $('form.addSpell').find('input[name="spellLvl"]').val('');
    });
}

function setSpellSet() {
    if (spells.length > 0) {
        if (filterLetter == '') {
            spellSet == spells;
        } else {
            spellSet = spells.filter(filterLetter);
            spellSet.sort(orderDoc);
        }
    }
}

function orderDoc(doc1, doc2) {
    var name1 = doc1.data.Name.toLowerCase();
    var name2 = doc2.data.Name.toLowerCase();
    if (name1 < name2) {
        return -1;
    }
    if (name1 > name2) {
        return 1;
    }
    return 0;
}

function filterLetter(doc) {
    var name = doc.data.Name.toLowerCase();
    return name.startsWith(letter);
}

function genLetters() {
    var m = {};
    var lets = [];
    for (var i = 0; i < spells.length; i++) {
        var l = spells[i].data.Name[0].toLowerCase();
        m[l] = '';
    }
    for (letter in m) {
        lets.push(letter);
    }
    lets.sort();
    if (lets.length > 0) {
        letter = lets[0];
    }
    for (var i = 0; i < lets.length; i++) {
        var s = $('<span class="label ' + ((lets[i] == letter) ? 'label-primary' : 'label-default') + ' letter" data-let="' + lets[i] + '">' + lets[i].toUpperCase() + '</span>');
        var o = $('<option value="' + lets[i] + '">' + lets[i].toUpperCase() + '</option>')
        $('div#letters').append(s);
        $('select#letters-dd').append(o);
    }
    if (cat == 'userC') {
        var s = $('<span class="label label-default letter" data-let="">All</span>');
        var o = $('<option value="">All</option>')
        $('div#letters').append(s);
        $('select#letters-dd').append(o);
    }
}

$(document).ready(function() {
    $('select#pageSize').change(function() {
        page = 0;
        pageSize =+ this.value;
        genResults(paginate(spellSet));
    });
});
