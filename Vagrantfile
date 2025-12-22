Vagrant.configure("2") do |config|
  config.vm.box =  "generic/ubuntu2204"

  config.vm.hostname = "student-api-prod"

  config.vm.synced_folder ".", "/vagrant"

  config.vm.synced_folder ".", "/vagrant",
    type: "rsync",
    rsync__auto: true

  config.vm.network "forwarded_port", guest: 8080, host: 8080

  config.vm.provider "virtualbox" do |vb|
    vb.memory = 2048
    vb.cpus = 2
  end

  config.vm.provision "shell", path: "provision/setup.sh"
end

