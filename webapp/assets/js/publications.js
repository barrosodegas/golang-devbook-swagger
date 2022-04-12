
$('#new-publication').on('submit', createPublication);

$(document).on('click', '.like-publication', likePublication);
$(document).on('click', '.unlike-publication', unlikePublication);

function createPublication(event) {
    event.preventDefault();

    $.ajax({
        url: '/publications',
        method: 'POST',
        data: {
            'title': $('#title').val(),
            'content': $('#content').val()
        }
    }).done(function() {
        window.location = '/home';
    }).fail(function(_) {
        Swal.fire('Ops...', 'Ocorreu um error ao tentar criar a publicação!', 'error');
    });
}

function likePublication(event) {
    event.preventDefault();

    const clickedElement = $(event.target);
    const publicationId = clickedElement.closest('div').data('publication-id');

    clickedElement.prop('disabled', true);
    
    $.ajax({
        url: `/publications/${publicationId}/like`,
        method: 'POST'
    }).done(function() {
        const likesCounter = clickedElement.next('span');
        const likes = parseInt(likesCounter.text());
        
        likesCounter.text(likes + 1);

        clickedElement.addClass("text-danger");
        clickedElement.addClass('unlike-publication');
        clickedElement.removeClass('like-publication');
    }).fail(function(_) {
        alert('Ocorreu um error ao tentar curtir a publicação!');
    }).always(function(_) {
        clickedElement.prop('disabled', false);
    });
}

function unlikePublication(event) {
    event.preventDefault();

    const clickedElement = $(event.target);
    const publicationId = clickedElement.closest('div').data('publication-id');

    clickedElement.prop('disabled', true);
    
    $.ajax({
        url: `/publications/${publicationId}/unlike`,
        method: 'POST'
    }).done(function() {
        const likesCounter = clickedElement.next('span');
        const likes = parseInt(likesCounter.text());
        
        likesCounter.text(likes - 1);

        clickedElement.addClass('like-publication');
        clickedElement.removeClass("text-danger");
        clickedElement.removeClass('unlike-publication');
    }).fail(function(_) {
        alert('Ocorreu um error ao tentar curtir a publicação!');
    }).always(function(_) {
        clickedElement.prop('disabled', false);
    });
}
