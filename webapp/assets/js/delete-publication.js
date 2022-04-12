
$('.delete-publication').on('click', deletePublication);

function deletePublication(event) {
    event.preventDefault();

    Swal.fire({
        title: 'Atenção',
        text: 'Tem certeza que deseja excluir esta publicação?',
        showCancelButton: true,
        cancelButtonText: 'Cancelar',
        icon: 'warning'
    }).then(function(confirmation) {

        if (!confirmation.value) return;

        const clickedElement = $(event.target);
        const publicationElement = clickedElement.closest('div');
        const publicationId = publicationElement.data('publication-id');
    
        clickedElement.prop('disabled', true);
    
        $.ajax({
            url: `/publications/${publicationId}`,
            method: 'DELETE'
        }).done(function() {
            publicationElement.fadeOut('slow', function () {
                $(this).remove();
            });
        }).fail(function(_) {
            Swal.fire('Ops...', 'Ocorreu um error ao tentar excluir a publicação!', 'error');
        }).always(function() {
            clickedElement.prop('disabled', false);
        });
    });
}
