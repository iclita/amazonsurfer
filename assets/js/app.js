$(document).ready(function(){
    console.log("Application started");
    $('#categories').select2({placeholder: "Select categories"});
});

function showSearchButton() {
    $('.searching').hide();
    $('#stop-button').hide();
    $('#search-button').show();
}

function showStopButton() {
    $('#search-button').hide();
    $('.searching').show();
    $('#stop-button').show();
}

function resetResultsTable() {
    $('#count').html(0);
    $('#count-text').show();
    $('#results tbody tr').remove();
    $('#results').hide();
}

function validateInput() {

    var isValid = true;

    if ($('#categories').val().length === 0) {
        isValid = false;
        $('.select2-selection').css('border', '2px solid red');
    } else {
        $('.select2-selection').css('border', '1px solid black');
    }

    var minPrice = $('#min-price').val() !== '' ? parseFloat($('#min-price').val()) : -1;
    var maxPrice = $('#max-price').val() !== '' ? parseFloat($('#max-price').val()) : -1;

    if (minPrice <= 0) {
        isValid = false;
        $('#min-price').addClass('error');
    } else {
        $('#min-price').removeClass('error');
    }

    if (maxPrice <= 0) {
        isValid = false;
        $('#max-price').addClass('error');
    } else {
        $('#max-price').removeClass('error');
    }

    var minBSR = $('#min-bsr').val() !== '' ? parseInt($('#min-bsr').val()) : -1;
    var maxBSR = $('#max-bsr').val() !== '' ? parseInt($('#max-bsr').val()) : -1;

    if (minBSR < 0) {
        isValid = false;
        $('#min-bsr').addClass('error');
    } else {
        $('#min-bsr').removeClass('error');
    }

    if (maxBSR <= 0) {
        isValid = false;
        $('#max-bsr').addClass('error');
    } else {
        $('#max-bsr').removeClass('error');
    }

    var minReviews = $('#min-reviews').val() !== '' ? parseInt($('#min-reviews').val()) : -1;
    var maxReviews = $('#max-reviews').val() !== '' ? parseInt($('#max-reviews').val()) : -1;

    if (minReviews < 0) {
        isValid = false;
        $('#min-reviews').addClass('error');
    } else {
        $('#min-reviews').removeClass('error');
    }

    if (maxReviews <= 0) {
        isValid = false;
        $('#max-reviews').addClass('error');
    } else {
        $('#max-reviews').removeClass('error');
    }

    var maxLength = $('#max-length').val() !== '' ? parseFloat($('#max-length').val()) : -1;

    if (maxLength <= 0) {
        isValid = false;
        $('#max-length').addClass('error');
    } else {
        $('#max-length').removeClass('error');
    }

    var maxWidth = $('#max-width').val() !== '' ? parseFloat($('#max-width').val()) : -1;

    if (maxWidth <= 0) {
        isValid = false;
        $('#max-width').addClass('error');
    } else {
        $('#max-width').removeClass('error');
    }

    var maxHeight = $('#max-height').val() !== '' ? parseFloat($('#max-height').val()) : -1;

    if (maxHeight <= 0) {
        isValid = false;
        $('#max-height').addClass('error');
    } else {
        $('#max-height').removeClass('error');
    }

    var maxWeight = $('#max-weight').val() !== '' ? parseFloat($('#max-weight').val()) : -1;

    if (maxWeight <= 0) {
        isValid = false;
        $('#max-weight').addClass('error');
    } else {
        $('#max-weight').removeClass('error');
    }

    var tolerance = $('#tolerance').val() !== '' ? parseFloat($('#tolerance').val()) : -1;

    if (tolerance < 0 || tolerance > 10) {
        isValid = false;
        $('#tolerance').addClass('error');
    } else {
        $('#tolerance').removeClass('error');
    }

    if (minPrice >= maxPrice) {
        isValid = false;
        $('#min-price').addClass('error');
        $('#max-price').addClass('error');
    } else {
        if (minPrice > 0 && maxPrice > 0) {
            $('#min-price').removeClass('error');
            $('#max-price').removeClass('error');
        }
    }

    if (minBSR >= maxBSR) {
        isValid = false;
        $('#min-bsr').addClass('error');
        $('#max-bsr').addClass('error');
    } else {
        if (minBSR >= 0 && maxBSR > 0) {
            $('#min-bsr').removeClass('error');
            $('#max-bsr').removeClass('error');
        }
    }

    if (minReviews >= maxReviews) {
        isValid = false;
        $('#min-reviews').addClass('error');
        $('#max-reviews').addClass('error');
    } else {
        if (minReviews >= 0 && maxReviews > 0) {
            $('#min-reviews').removeClass('error');
            $('#max-reviews').removeClass('error');
        }
    }

    return isValid;
}