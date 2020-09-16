/**
 * Copyright 2020 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

const {Datastore} = require('@google-cloud/datastore');

const datastore = new Datastore();

exports.key = datastore.key;

exports.get = (key) => {
	return datastore.get(key).then(entity => {
		if (entity.length == 0) {
			return undefined;
		}
		return entity[0];
	});
};

exports.set = (key, entity) => {
	return datastore.insert({key: key, data: entity});
};

exports.merge = (key, entity) => {
	return datastore.merge({key: key, data: entity});
};
