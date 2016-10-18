require 'active_record'

def setup
	ActiveRecord::Base.establish_connection(
		adapter: 'sqlite3',
		database: '../../tmp/db.sqlite3',
	)
end

def define_schema
	ActiveRecord::Schema.define do
		create_table :builds do |t|
			t.string :log
			t.string :status
			t.binary :source_file
		end
	end
end

class Build < ActiveRecord::Base
end

# entry point
ActiveRecord::Base.logger = Logger.new(STDERR)
setup
con = ActiveRecord::Base.connection
define_schema unless con.data_source_exists?('builds')

# ./app
system('mkdir -p tmp; cd app; zip -r -v ../tmp/app.zip ./')
unless Build.where(id: 10000).exists?
  Build.create id: 10000, source_file: File.read('tmp/app.zip')
end

puts 'finished :-)'
