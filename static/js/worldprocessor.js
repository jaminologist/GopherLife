$(document).ready(function(){

    getWorld()

    $(document).keydown(function(e) {
        var key = e.which;

        var leftArrow = 37
        var rightArrow = 39
        var upArrow = 38
        var downArrow = 40

        var inputData = new FormData()
        inputData.append("keydown", key)



        $.ajax({
            type: 'GET',
            url: '/ShiftWorldView?keydown=' + key,
            error: function(xhr, statusText, err) {
                console.log("ERR")
              alert("error"+xhr.status);
            },
            success: function(){
                console.log("syccc")
            }
         })
    });
    

})


function getWorld(){


   /* var xhr = new XMLHttpRequest();

    xhr.open("GET", "/PollWorld")
    xhr.send(null)

    xhr.onreadystatechange() = function(event) {

        if(xhr.readyState === 4){
            $("#worldDiv").html(data)
        }

    }*/

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