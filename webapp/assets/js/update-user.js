
$('#form-update-user').on('submit', updateUser);

function updateUser(event) {
    event.preventDefault();

    $('#btn-update-user').prop('disabled', true);

    $.ajax({
        url: '/me',
        method: 'PUT',
        data: {
            'name': $('#name').val(),
            'nick': $('#nick').val(),
            'email': $('#email').val()
        }
    }).done(function() {
        Swal.fire(
            'Sucesso!',
            'Usuário atualizada com sucesso!',
            'success'
        ).then(function() {
            window.location = '/me';
        });
    }).fail(function(_) {
        Swal.fire('Ops...', 'Ocorreu um error ao tentar atualizar seus dados!', 'error');
    }).always(function() {
        $('#btn-update-user').prop('disabled', false);
    });
}
