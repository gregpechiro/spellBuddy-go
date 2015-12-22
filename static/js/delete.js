$(document).ready(function() {

	$('.delete-button').click(function() {
		$('.alert-dismissable').alert('close');
    	$('form#deleteForm').attr('action', $(this).attr('data-delete'));
    	$('span#message').html($(this).attr('data-message'));
    	$('div#deleteAlert').removeClass('hide');
    });

    $('a#deleteCancel').click(function() {
    	$('form#deleteForm').attr('action', '');
		$('span#message').html('');
    	$('div#deleteAlert').addClass('hide');
    });

});
