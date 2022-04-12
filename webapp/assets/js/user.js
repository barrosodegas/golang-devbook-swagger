
$('#unfollow').on('click', unfollow);
$('#follow').on('click', follow);
$('#delete-user').on('click', deleteUser);

function unfollow(event) {
    event.preventDefault();

    const userId = $(this).data('user-id');
    $(this).prop('disabled', true);

    $.ajax({
        url: `/users/${userId}/unfollow`,
        method: 'POST'
    }).done(function() {
        window.location = `/users/${userId}`
    }).fail(function(_) {
        Swal.fire('Ops...', 'Ocorreu um erro ao tentar deixar de seguir o usuário!', 'error');
        $(this).prop('disabled', false);
    });
}

function follow(event) {
    event.preventDefault();

    const userId = $(this).data('user-id');
    $(this).prop('disabled', true);

    $.ajax({
        url: `/users/${userId}/follow`,
        method: 'POST'
    }).done(function() {
        window.location = `/users/${userId}`
    }).fail(function(_) {
        Swal.fire('Ops...', 'Ocorreu um erro ao tentar seguir o usuário!', 'error');
        $(this).prop('disabled', false);
    });
}

function deleteUser(event) {
    event.preventDefault();

    Swal.fire({
        title: 'Atenção',
        text: 'Tem certeza que deseja excluir a sua conta?',
        showCancelButton: true,
        cancelButtonText: 'Cancelar',
        icon: 'warning'
    }).then(function(confirmation) {

        if (!confirmation.value) return;

        $.ajax({
            url: `/users`,
            method: 'DELETE'
        }).done(function() {
            Swal.fire(
                'Sucesso!',
                'Conta excluída com sucesso!',
                'success'
            ).then(function() {
                window.location = '/logout';
            });
        }).fail(function(_) {
            Swal.fire('Ops...', 'Ocorreu um erro ao tentar excluir o usuário!', 'error');
            $(this).prop('disabled', false);
        });
    });
}