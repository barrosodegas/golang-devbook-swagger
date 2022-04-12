
$('#form-create-user').on('submit', createUser);

function createUser(event) {
    event.preventDefault();

    var password = $('#password').val();
    var confirmPassword = $('#confirm-password').val();
    
    if (password != confirmPassword) {
        Swal.fire('Ops...', 'As senhas não coincidem!', 'error');
        return;
    }

    $.ajax({
        url: '/users',
        method: 'POST',
        data: {
            'name': $('#name').val(),
            'nick': $('#nick').val(),
            'email': $('#email').val(),
            'password': $('#password').val()
        }
    }).done(function() {
        Swal.fire('Sucesso!', 'Usuário criado com sucesso!', 'success')
            .then(function() {
                $.ajax({
                    url: '/login',
                    method: 'POST',
                    data: {
                        'email': $('#email').val(),
                        'password': $('#password').val()
                    }
                }).done(function() {
                    window.location = '/home';
                }).fail(function(_) {
                    Swal.fire('Ops...', 'Erro ao tentar realizar o login! Tente novamente mais tarde,', 'error');
                });
            });
    }).fail(function(_) {
        Swal.fire('Ops...', 'Erro ao tentar cadastrar o usuário! Tente novamente mais tarde.', 'error');
    });
}