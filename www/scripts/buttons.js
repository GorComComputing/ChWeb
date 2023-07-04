		$(document).ready(function(){
		
			$('#formSave').submit(function(){
				$.ajax({
					type: "POST",
					url: "/save",
					data: "Leapsectz="+$("#Leapsectz").val() +"&"+
					"Driftfile="+$("#Driftfile").val() +"&"+
					"Makestep="+$("#Makestep").val() +"&"+
					"Makestep2="+$("#Makestep2").val() +"&"+
					"Rtcsync="+$("#Rtcsync").val() +"&"+
					"Logdir="+$("#Logdir").val() +"&"+
					"LocalStratum="+$("#LocalStratum").val() +"&"+
					"Server="+$("#Server").val() +"&"+
					"Allow="+$("#Allow").val(),
					success: function(html){
						//$("#usage").html(html);
    						$('textarea').val(html);
				   	}
				});
				return false;
			});
			
			$('#formStart').submit(function(){
				$.ajax({
					type: "POST",
					url: "/start",
					//success: function(html){
						//$("#msg").html(html);
				   	//}
				});
				return false;
			});
			
			$('#formStop').submit(function(){
				$.ajax({
					type: "POST",
					url: "/stop",
					//success: function(html){
						//$("#msg").html(html);
				   	//}
				});
				return false;
			});
			
			$('#formRestart').submit(function(){
				$.ajax({
					type: "POST",
					url: "/restart",
					//success: function(html){
						//$("#msg").html(html);
				   	//}
				});
				return false;
			});
			
			$('#formActivity').submit(function(){
				$.ajax({
					type: "POST",
					url: "/activity",
					success: function(html){
						//$("#usage").html(html);
    						$('textarea').val(html);
				   	}
				});
				return false;
			});
			
			$('#formTracking').submit(function(){
				$.ajax({
					type: "POST",
					url: "/tracking",
					success: function(html){
						//$("#usage").html(html);
    						$('textarea').val(html);
				   	}
				});
				return false;
			});
			
			$('#formSources').submit(function(){
				$.ajax({
					type: "POST",
					url: "/sources",
					success: function(html){
						//$("#usage").html(html);
    						$('textarea').val(html);
				   	}
				});
				return false;
			});
			
			$('#formSourceStats').submit(function(){
				$.ajax({
					type: "POST",
					url: "/sourcestats",
					success: function(html){
						//$("#usage").html(html);
    						$('textarea').val(html);
				   	}
				});
				return false;
			});
			
			$('#formClients').submit(function(){
				$.ajax({
					type: "POST",
					url: "/clients",
					success: function(html){
						//$("#usage").html(html);
    						$('textarea').val(html);
				   	}
				});
				return false;
			});
			
			$('#formConfig').submit(function(){
				$.ajax({
					type: "POST",
					url: "/config",
					success: function(html){
						//$("#usage").html(html);
    						$('textarea').val(html);
				   	}
				});
				return false;
			});
			
			
			
			$('#formSaveConfig').submit(function(){
				$.ajax({
					type: "POST",
					url: "/saveconfig",
					data: "file="+$("#usage").val(),
					/*success: function(html){
						//$("#usage").html(html);
    						$('textarea').val(html);
				   	}*/
				});
				return false;
			});
			
		});
