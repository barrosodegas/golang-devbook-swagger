
$('#form-update-password').on('submit', updatePassword);

function updatePassword(event) {
    event.preventDefault();

    if($('#new-password').val() != $('#confirm-new-password').val()) {
        Swal.fire('Ops...', 'As senhas não coincidem!', 'warning');
        return;
    }

    $('#btn-update-password').prop('disabled', true);

    $.ajax({
        url: '/me/update-password',
        method: 'POST',
        data: {
            'current': $('#password').val(),
            'new': $('#new-password').val()
        }
    }).done(function() {
        Swal.fire(
            'Sucesso!',
            'Senha atualizada com sucesso!',
            'success'
        ).then(function() {
            window.location = '/me';
        });
    }).fail(function(_) {
        Swal.fire('Ops...', 'Ocorreu um error ao tentar atualizar sua senha!', 'error');
        $('#btn-update-password').prop('disabled', false);
    });
}
