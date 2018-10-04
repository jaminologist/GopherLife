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


    $(document).on("click", ".gopher", function(){

        var position = $(this).attr("id")

        $.ajax({
            type: 'GET',
            url: '/SelectGopher?position=' + position,
            error: function(xhr, statusText, err) {
                console.log("ERR")
              alert("error"+xhr.status);
            },
            success: function(data){
                $("#worldDiv").html(data)
            }
         })

    });


})


function getWorld(){

    $.ajax({
        url: '/PollWorld',
        success: function(data) {
              $("#worldDiv").html(data)
        },

        type: 'GET'
     }).always(function(){
        setTimeout(getWorld(),5000);
     });

}
