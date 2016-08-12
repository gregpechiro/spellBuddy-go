function renderPicked() {
    $('#spellAccordion').html('');
    for (var i = 0; i < picked.length; i ++) {
        if (picked[i].length > 0) {
            var panel = $('<div class="panel panel-default">' +
                '<div class="panel-heading" role="tab" id="heading-' + i + '">' +
                    '<h4 class="panel-title">' +
                        '<a data-toggle="collapse" data-parent="#spellAccordion" style="cursor:pointer;" data-target="#collapse-' + i + '">' +
                            'Level ' + i + ' Spells'+
                        '</a>' +
                    '</h4>' +
                '</div>' +
                '<div id="collapse-' + i + '" class="panel-collapse collapse">' +
                    '<div id="body-' + i + '" class="panel-body">' +
                    '</div>' +
                '</div>' +
            '</div>');
            for (var j = 0; j < picked[i].length; j++) {
                var spell = $('<li class="list-group-item clearfix picked-li">' +
                    picked[i][j] +
                    '<span class="visible-xs-block"><br></span>'+
                    '<div class="input-group pull-right col-lg-6 input-group-sm">' +
                        '<input id="change-' + i + '-' + j + '" type="number" class="form-control" placeholder="Level">' +
                        '<span class="input-group-btn">' +
                            '<button id="change-' + i + '-' + j + '" class="btn btn-default change-lvl" data-lvl="' + i + '" data-idx="' + j + '">Change Level</button>' +
                            '<button id="delete-spell" class="btn btn-danger" data-lvl="' + i + '" data-idx="' + j + '">Delete</button>' +
                        '</span>' +
                    '</div>' +
                '</li>');
                panel.find('#body-' + i).append(spell);
            }
            $('#spellAccordion').append(panel)
        }
    }
    $('button#delete-spell').click(function() {
        var dat = 'userId=' + userId + '&spellLvl=' + $(this).attr('data-lvl') + '&idx=' + $(this).attr('data-idx');
        $.ajax({
            type: 'POST',
            url: '/user/delSpell',
            data: dat,
            success: function(resp) {
                if (!resp.success) {
                    showError(resp.msg);
                    return
                }
                picked = resp.picked;
                renderPicked();
            },
            error: function() {
                showError("Error deleting spell");
            }
        });
    });
    $('button.change-lvl').click(function() {
        var newLvlStr = $('input#' + this.id).val();
        var newLvl =+ $('input#' + this.id).val();
        if (newLvl < 0) {
            showError("Spell level cannot be negative");
            $('input#' + this.id).val('');
            return
        }
        if (newLvl > 9) {
            showError("Spell level must be between 0 and 9");
            $('input#' + this.id).val('');
            return
        }
        if (newLvlStr === '') {
            showError("Spell level cannot be empty");
            $('input#' + this.id).val('');
            return
        }

        var dat = 'userId=' + userId + '&spellLvl=' + $(this).attr('data-lvl') + '&idx=' + $(this).attr('data-idx') + '&newLvl=' + newLvl;

        $.ajax({
            type: 'POST',
            url: '/user/changeLvl',
            data: dat,
            success: function(resp) {
                if (!resp.success) {
                    showError(resp.msg);
                    return
                }
                picked = resp.picked;
                renderPicked();
            },
            error: function() {
                showError("Error changing spell level");
            }
        });

    });
}

function orderName(name1, name2) {
    name1 = name1.toLowerCase();
    name2 = name2.toLowerCase();
    if (name1 < name2) {
        return -1;
    }
    if (name1 > name2) {
        return 1;
    }
    return 0;
}
