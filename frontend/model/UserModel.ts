export class UserModel {
  name: string;
  email: string;
  password: string;
  age: number;
  bio: string;
  userType: string; // VOLUNTEER, ELDERLY, INSTITUTION

  constructor(
    name: string,
    email: string,
    password: string,
    age: number,
    bio: string,
    userType: string,
  ) {
    this.name = name;
    this.email = email;
    this.password = password;
    this.age = age;
    this.bio = bio;
    this.userType = userType;
  }
}