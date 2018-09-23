$(document).ready(function(){

    getWorld()
    

})


function getWorld(){

    $.ajax({
        url: '/PollWorld',
        error: function(xhr, statusText, err) {
          //alert("error"+xhr.status);
        },

        success: function(data) {
            //alert(data)
            $("#worldDiv").html(data)
        },
        
        type: 'GET'
     }).always(function(){
        setTimeout(getWorld(),5000);
     });

}