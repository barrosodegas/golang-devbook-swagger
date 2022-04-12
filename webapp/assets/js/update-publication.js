
$('#edit-publication').on('submit', updatePublication);

function updatePublication(event) {
    event.preventDefault();

    $('#btn-update-pub').prop('disabled', true);

    var publicationId = $('#btn-update-pub').data('publication-id');

    $.ajax({
        url: `/publications/${publicationId}`,
        method: 'PUT',
        data: {
            'title': $('#title').val(),
            'content': $('#content').val()
        }
    }).done(function() {
        Swal.fire(
            'Sucesso!',
            'Publicação atualizada com sucesso!',
            'success'
        ).then(function() {
            window.location = '/home';
        });
    }).fail(function(_) {
        Swal.fire('Ops...', 'Ocorreu um error ao tentar atualizar a publicação!', 'error');
    }).always(function() {
        $('#btn-update-pub').prop('disabled', false);
    });
}
