using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace Web.Migrations
{
    /// <inheritdoc />
    public partial class MachineClientFeat : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<Guid>(
                name: "MachineClientId",
                table: "Tokens",
                type: "uuid",
                nullable: false,
                defaultValue: new Guid("00000000-0000-0000-0000-000000000000"));

            migrationBuilder.CreateTable(
                name: "MachineClients",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    ClientId = table.Column<string>(type: "text", nullable: false),
                    ClientSecretHash = table.Column<string>(type: "text", nullable: false),
                    Description = table.Column<string>(type: "text", nullable: true),
                    CreatedDate = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    LastModifiedDate = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_MachineClients", x => x.Id);
                });

            migrationBuilder.CreateIndex(
                name: "IX_Tokens_MachineClientId",
                table: "Tokens",
                column: "MachineClientId");

            migrationBuilder.AddForeignKey(
                name: "FK_Tokens_MachineClients_MachineClientId",
                table: "Tokens",
                column: "MachineClientId",
                principalTable: "MachineClients",
                principalColumn: "Id",
                onDelete: ReferentialAction.Cascade);
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_Tokens_MachineClients_MachineClientId",
                table: "Tokens");

            migrationBuilder.DropTable(
                name: "MachineClients");

            migrationBuilder.DropIndex(
                name: "IX_Tokens_MachineClientId",
                table: "Tokens");

            migrationBuilder.DropColumn(
                name: "MachineClientId",
                table: "Tokens");
        }
    }
}
